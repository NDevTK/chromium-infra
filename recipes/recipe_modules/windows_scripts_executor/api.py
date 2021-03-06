# Copyright 2021 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from recipe_engine import recipe_api
from recipe_engine.recipe_api import Property
# Windows command helpers
from . import offline_winpe_customization as offwinpecust
from . import sources
from . import helper
from PB.recipes.infra.windows_image_builder import windows_image_builder as wib
from PB.recipes.infra.windows_image_builder import sources as src_pb


class WindowsPSExecutorAPI(recipe_api.RecipeApi):
  """API for using Windows PowerShell scripts."""

  def __init__(self, *args, **kwargs):
    super(WindowsPSExecutorAPI, self).__init__(*args, **kwargs)
    self._scripts = self.resource('WindowsPowerShell\Scripts')
    self._workdir = ''
    self._sources = None
    self._configs_dir = None
    self._customizations = []

  def init(self, config):
    """ init initializes all the dirs and sub modules required.
        Args:
          config: wib.Image proto object representing the image to be created
    """
    # ensure that arch is specified in the image
    if config.arch == wib.Arch.ARCH_UNSPECIFIED:
      raise self.m.step.StepFailure('Missing arch in config')
    arch = wib.Arch.Name(config.arch).replace('ARCH_', '').lower()

    self._sources = sources.Source(self.m.path['cache'].join('Pkgs'),
                                   self.m.step, self.m.path, self.m.file,
                                   self.m.raw_io, self.m.cipd, self.m.gsutil,
                                   self.m.gitiles, self.m.git)

    self._configs_dir = self.m.path['cleanup'].join('configs')
    helper.ensure_dirs(self.m.file, [self._configs_dir])

    # initialize all customizations
    for cust in config.customizations:
      if cust.WhichOneof('customization') == 'offline_winpe_customization':
        self._customizations.append(
            offwinpecust.OfflineWinPECustomization(
                cust,
                arch=arch,
                scripts=self._scripts,
                configs=self._configs_dir,
                step=self.m.step,
                path=self.m.path,
                powershell=self.m.powershell,
                m_file=self.m.file,
                source=self._sources))

  def pin_available_sources(self):
    """ pin_wib_config pins the given config to current refs."""
    with self.m.step.nest('Pin all the required artifacts'):
      self._sources.pin()

  def gen_canonical_configs(self, config):
    """ gen_canonical_configs strips all the names in the config and returns
        individual configs containing one customization per image.
        Args:
          config: wib.Image proto representing the image to be generated
        Example:
          Given an Image
            Image{
              arch: x86,
              name: "windows10_x86_GCE",
              customizations: [
                Customization{
                  OfflineWinPECustomization{
                    name: "winpe_networking"
                    image_dest: GCSSrc{
                      bucket: "chrome-win-wim"
                      source: "rel/win10_networking.wim"
                    }
                    ...
                  }
                },
                Customization{
                  OfflineWinPECustomization{
                    name: "winpe_diskpart"
                    image_src: Src{
                      gcs_src: GCSSrc{
                        bucket: "chrome-win-wim"
                        source: "rel/win10_networking.wim"
                      }
                    }
                    ...
                  }
                }
              ]
            }
          Writes two configs: windows10_x86_GCE-winpe_networking.cfg with
            Image{
              arch: x86,
              name: "",
              customizations: [
                Customization{
                  OfflineWinPECustomization{
                    name: ""
                    image_dest: GCSSrc{
                      bucket: "chrome-win-wim"
                      source: "rel/win10_networking.wim"
                    }
                    ...
                  }
               }
              ]
            }
          and windows10_x86_GCE-winpe_diskpart.cfg with
            Image{
              arch: x86,
              name: "",
              customizations: [
                Customization{
                  OfflineWinPECustomization{
                    name: ""
                    image_src: Src{
                      gcs_src: GCSSrc{
                        bucket: "chrome-win-wim"
                        source: "rel/win10_networking.wim"
                      }
                    }
                    ...
                  }
                }
              ]
            }
          to disk, calculates the hash for each config and sets the key for each
          of them. The strings representing name of the image, customization,...
          etc,. are set to empty before calculating the hash to maintain the
          uniqueness of the hash.
    """
    for cust in self._customizations:
      # create a new image object, with same arch and containing only one
      # customization
      canon_image = wib.Image(
          arch=config.arch, customizations=[cust.get_canonical_cfg()])
      name = cust.name()
      # write the config to disk
      cfg_file = self._configs_dir.join('{}-{}.cfg'.format(config.name, name))
      self.m.file.write_proto(
          'Write config {}'.format(cfg_file),
          cfg_file,
          canon_image,
          codec='TEXTPB')
      # estimate the unique hash for the config (identifier for the image built
      # by this config)
      key = self.m.file.file_hash(cfg_file)
      cust.set_key(key)
      # save the config to disk as <key>.cfg
      key_file = self._configs_dir.join('{}.cfg'.format(key))
      self.m.file.copy('Copy {} to {}'.format(cfg_file, key_file), cfg_file,
                       key_file)

  def download_available_packages(self):
    """ download_available_packages downloads the src refs that are pinned """
    with self.m.step.nest('Download all available packages'):
      self._sources.download()

  def execute_config(self, config):
    """ Executes the windows image builder user config.
        Args:
          config: wib.Image proto representing the image to be generated
    """
    with self.m.step.nest('execute config {}'.format(config.name)):
      for cust in self._customizations:
        output = cust.get_output()
        if not self._sources.exists(src_pb.Src(gcs_src=output.gcs_src)):
          # execute the customization if we don't have the output
          cust.execute_customization()

  def upload_wib_artifacts(self):
    """ upload_wib_artifacts uploads all the available artifacts """
    self._sources.upload()
