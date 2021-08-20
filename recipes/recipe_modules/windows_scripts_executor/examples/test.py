# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from PB.recipes.infra.windows_image_builder import windows_image_builder as wib

from recipe_engine.post_process import DropExpectation, StatusFailure
from recipe_engine.post_process import StatusSuccess, StepCommandRE

DEPS = [
    'depot_tools/gitiles',
    'windows_scripts_executor',
    'recipe_engine/properties',
    'recipe_engine/platform',
    'recipe_engine/json',
]

PROPERTIES = wib.Image


def RunSteps(api, image):
  api.windows_scripts_executor.pin_wib_config(image)
  api.windows_scripts_executor.download_wib_artifacts(image)
  api.windows_scripts_executor.execute_wib_config(image)


def GenTests(api):
  # various step data for testing
  GEN_WPE_MEDIA_FAIL = api.step_data(
      'execute config win10_2013_x64.offline winpe ' +
      'customization offline_winpe_2013_x64.generate ' +
      'windows image folder for x86 in ' +
      '[CACHE]\\WinPEImage.PowerShell> Gen WinPE media for x86',
      stdout=api.json.output({
          'results': {
              'Success': False,
              'Command': 'powershell',
              'ErrorInfo': {
                  'Message': 'Failed step'
              },
          }
      }))

  GEN_WPE_MEDIA_PASS = api.step_data(
      'execute config win10_2013_x64.offline winpe ' +
      'customization offline_winpe_2013_x64.generate ' +
      'windows image folder for x86 in ' +
      '[CACHE]\\WinPEImage.PowerShell> Gen WinPE media for x86',
      stdout=api.json.output({'results': {
          'Success': True,
      }}))

  MOUNT_WIM_PASS = api.step_data(
      'execute config win10_2013_x64.offline winpe customization ' +
      'offline_winpe_2013_x64.generate windows image folder for ' +
      'x86 in [CACHE]\\WinPEImage.PowerShell> Mount wim to ' +
      '[CACHE]\\WinPEImage\\mount',
      stdout=api.json.output({
          'results': {
              'Success': True
          },
      }))

  UMOUNT_WIM_PASS = api.step_data(
      'execute config win10_2013_x64.offline winpe ' +
      'customization offline_winpe_2013_x64.PowerShell> ' +
      'Unmount wim at [CACHE]\\WinPEImage\\mount',
      stdout=api.json.output({
          'results': {
              'Success': True
          },
      }))

  PIN_FILE_STARTNET_PASS = api.step_data(
      'Pin git artifacts to refs.gitiles log: ' +
      'HEAD/windows/artifacts/startnet.cmd',
      api.gitiles.make_log_test_data('HEAD'),
  )

  FETCH_FILE_STARTNET_PASS = api.step_data(
      'Get all git artifacts.fetch ' +
      'ef70cb069518e6dc3ff24bfae7f195de5099c377:' +
      'windows/artifacts/startnet.cmd',
      api.gitiles.make_encoded_file('Wpeinit'))

  ADD_FILE_STARTNET_PASS = api.step_data(
      'execute config win10_2013_x64.offline winpe customization ' +
      'offline_winpe_2013_x64.PowerShell> Add file ' +
      '[CACHE]\\GITPkgs\\ef70cb069518e6dc3ff24bfae7f195de5099c377\\' +
      'windows\\artifacts\\startnet.cmd',
      stdout=api.json.output({'results': {
          'Success': True,
      }}))

  ADD_FILE_STARTNET_FAIL = api.step_data(
      'execute config win10_2013_x64.offline winpe customization ' +
      'offline_winpe_2013_x64.PowerShell> Add file ' +
      '[CACHE]\\GITPkgs\\ef70cb069518e6dc3ff24bfae7f195de5099c377\\' +
      'windows\\artifacts\\startnet.cmd',
      stdout=api.json.output({
          'results': {
              'Success': False,
              'Command': 'powershell',
              'ErrorInfo': {
                  'Message': 'Failed step',
              },
          }
      }))

  ADD_FILE_CIPD_PASS = api.step_data(
      'execute config win10_2013_x64.offline ' +
      'winpe customization offline_winpe_2013_x64.PowerShell> ' +
      'Add file [CACHE]\\CIPDPkgs\\resolved-instance_id-of-latest----------' +
      '\\infra_internal\\labs\\drivers\\microsoft\\windows_adk\\winpe' +
      '\\winpe-dot3svc\\windows-amd64\\*',
      stdout=api.json.output({'results': {
          'Success': True,
      }}))

  # actions for adding files
  ACTION_ADD_STARTNET = wib.Action(
      add_file=wib.AddFile(
          name='add_startnet_file',
          src=wib.Src(
              git_src=wib.GITSrc(
                  repo='chromium.dev',
                  ref='HEAD',
                  src='windows/artifacts/startnet.cmd'),),
          dst='Windows\\System32',
      ))

  ACTION_ADD_DOT3SVC = wib.Action(
      add_file=wib.AddFile(
          name='add winpe-dot3svc',
          src = wib.Src(
              cipd_src=wib.CIPDSrc(
              package='infra_internal/labs/drivers/' +
              'microsoft/windows_adk/winpe/' + 'winpe-dot3svc',
              refs='latest',
              platform='windows-amd64',
              ),
          ),
          dst='Windows\\System32\\',
      ))

  # Post process check for save and discard options during unmount
  UMOUNT_PP_DISCARD = api.post_process(
      StepCommandRE, 'execute config win10_2013_x64.offline winpe ' +
      'customization offline_winpe_2013_x64.PowerShell> ' +
      'Unmount wim at [CACHE]\\WinPEImage\\mount', [
          'python', '-u',
          'RECIPE_MODULE\[infra::powershell\]\\\\resources\\\\psinvoke.py',
          '--command', 'Dismount-WindowsImage', '--logs',
          '\[CLEANUP\]\\\\Dismount-WindowsImage', '--',
          '-Path "\[CACHE\]\\\\WinPEImage\\\\mount"',
          '-LogPath "\[CLEANUP\]\\\\Dismount-WindowsImage\\\\unmount.log"',
          '-LogLevel WarningsInfo', '-Discard'
      ])

  UMOUNT_PP_SAVE = api.post_process(
      StepCommandRE, 'execute config win10_2013_x64.offline winpe ' +
      'customization offline_winpe_2013_x64.PowerShell> ' +
      'Unmount wim at [CACHE]\\WinPEImage\\mount', [
          'python', '-u',
          'RECIPE_MODULE\[infra::powershell\]\\\\resources\\\\\psinvoke.py',
          '--command', 'Dismount-WindowsImage', '--logs',
          '\[CLEANUP\]\\\\Dismount-WindowsImage', '--',
          '-Path "\[CACHE\]\\\\WinPEImage\\\\mount"',
          '-LogPath "\[CLEANUP\]\\\\Dismount-WindowsImage\\\\unmount.log"',
          '-LogLevel WarningsInfo', '-Save'
      ])

  yield (
      api.test('Fail win image folder creation', api.platform('win', 64)) +
      api.properties(
          wib.Image(
              name='win10_2013_x64',
              arch=wib.ARCH_X86,
              offline_winpe_customization=wib.OfflineCustomization(
                  name='offline_winpe_2013_x64',
                  offline_customization=[
                      wib.OfflineAction(
                          name='network_setup', actions=[ACTION_ADD_STARTNET])
                  ]))) +
      PIN_FILE_STARTNET_PASS +  # pin the startnet file to current refs
      FETCH_FILE_STARTNET_PASS +  # fetch the startnet file from gitiles
      GEN_WPE_MEDIA_FAIL +  # Fail to create a winpe media folder
      api.post_process(StatusFailure) +  # recipe should fail
      api.post_process(DropExpectation))

  yield (api.test('Missing arch in config', api.platform('win', 64)) +
         api.properties(
             wib.Image(
                 name='win10_2013_x64',
                 offline_winpe_customization=wib.OfflineCustomization(
                     name='offline_winpe_2013_x64',))) +
         api.post_process(StatusFailure) +  # recipe should fail
         api.post_process(DropExpectation))

  yield (
      api.test('Fail add file step', api.platform('win', 64)) + api.properties(
          wib.Image(
              name='win10_2013_x64',
              arch=wib.ARCH_X86,
              offline_winpe_customization=wib.OfflineCustomization(
                  name='offline_winpe_2013_x64',
                  offline_customization=[
                      wib.OfflineAction(
                          name='network_setup', actions=[ACTION_ADD_STARTNET])
                  ]))) + GEN_WPE_MEDIA_PASS + MOUNT_WIM_PASS +
      PIN_FILE_STARTNET_PASS +  # pin the git file to current refs
      FETCH_FILE_STARTNET_PASS +  # fetch the file from gitiles
      ADD_FILE_STARTNET_FAIL +  # Fail to add file
      UMOUNT_WIM_PASS +  # Unmount the wim
      UMOUNT_PP_DISCARD +  # Discard the changes made to wim
      api.post_process(StatusFailure) +  # recipe fails
      api.post_process(DropExpectation))

  yield (
      api.test('Add file from cipd', api.platform('win', 64)) + api.properties(
          wib.Image(
              name='win10_2013_x64',
              arch=wib.ARCH_X86,
              offline_winpe_customization=wib.OfflineCustomization(
                  name='offline_winpe_2013_x64',
                  offline_customization=[
                      wib.OfflineAction(
                          name='network_setup', actions=[
                              ACTION_ADD_DOT3SVC,
                          ])
                  ]))) + GEN_WPE_MEDIA_PASS +  # generate the winpe media
      MOUNT_WIM_PASS +  # mount the generated wim
      ADD_FILE_CIPD_PASS +  # add the file from cipd
      UMOUNT_WIM_PASS +  # Unmount the wim
      UMOUNT_PP_SAVE +  # Check if the changes are saved to wim
      api.post_process(StatusSuccess) + api.post_process(DropExpectation))

  yield (
      api.test('Add file from git', api.platform('win', 64)) + api.properties(
          wib.Image(
              name='win10_2013_x64',
              arch=wib.ARCH_X86,
              offline_winpe_customization=wib.OfflineCustomization(
                  name='offline_winpe_2013_x64',
                  offline_customization=[
                      wib.OfflineAction(
                          name='network_setup', actions=[
                              ACTION_ADD_STARTNET,
                          ])
                  ]))) + PIN_FILE_STARTNET_PASS +  # pin the startnet refs
      FETCH_FILE_STARTNET_PASS +  # fetch the startnet file
      GEN_WPE_MEDIA_PASS +  # successfully gen winpe media
      MOUNT_WIM_PASS +  # mount the wim
      ADD_FILE_STARTNET_PASS +  # Add the downloaded file
      UMOUNT_WIM_PASS +  # Unmount the wim
      UMOUNT_PP_SAVE +  # Check unmount didn't discard the changes
      api.post_process(StatusSuccess) + api.post_process(DropExpectation))

  yield (api.test('Happy path', api.platform('win', 64)) + api.properties(
      wib.Image(
          name='win10_2013_x64',
          arch=wib.ARCH_X86,
          offline_winpe_customization=wib.OfflineCustomization(
              name='offline_winpe_2013_x64',
              offline_customization=[
                  wib.OfflineAction(
                      name='network_setup',
                      actions=[
                          ACTION_ADD_STARTNET,
                          ACTION_ADD_DOT3SVC,
                      ])
              ]))) + GEN_WPE_MEDIA_PASS +  # generate the winpe media
         MOUNT_WIM_PASS +  # mount the generated wim
         PIN_FILE_STARTNET_PASS +  # pin the startnet file to current refs
         FETCH_FILE_STARTNET_PASS +  # fetch the startnet file
         ADD_FILE_STARTNET_PASS +  # add file from git to wim
         ADD_FILE_CIPD_PASS +  # add the file from cipd to wim
         UMOUNT_WIM_PASS +  # unmount the wim
         UMOUNT_PP_SAVE +  # Save the changes made to the wim
         api.post_process(StatusSuccess) + api.post_process(DropExpectation))
