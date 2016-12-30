# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from crash import callstack_detectors
from crash.callstack_detectors import StartOfCallStack
from crash.flag_manager import ParsingFlag
from crash.flag_manager import FlagManager
from crash.test.stacktrace_test_suite import StacktraceTestSuite
from crash.type_enums import CallStackFormatType
from crash.type_enums import LanguageType


class CallStackDetectorTest(StacktraceTestSuite):

  def testAndroidJobDetector(self):
    """Tests that ``AndroidJobDetector`` detects android job callstack."""
    stack_detector = callstack_detectors.AndroidJobDetector()
    flag_manager = FlagManager()
    flag_manager.Register('group',
                          ParsingFlag('java_main_stack_flag', value=True))

    self.assertTupleEqual(
        stack_detector('java.lang.IllegalStateException: blabla', flag_manager),
        StartOfCallStack(0, CallStackFormatType.JAVA, LanguageType.JAVA, {}))
    self.assertTupleEqual(
        stack_detector('org.chromium.src.BlaBla', flag_manager),
        StartOfCallStack(1, CallStackFormatType.JAVA, LanguageType.JAVA, {}))
    self.assertTupleEqual(
        stack_detector('Caused by:', flag_manager),
        StartOfCallStack(1, CallStackFormatType.JAVA, LanguageType.JAVA, {}))
    self.assertTupleEqual(
        stack_detector('com.google.android.BlaBla', flag_manager),
        StartOfCallStack(1, CallStackFormatType.JAVA, LanguageType.JAVA, {}))
    self.assertIsNone(stack_detector('dummy', flag_manager))

  def testSyzyasanDetector(self):
    """Tests that ``SyzyasanDetector`` detects sysyasn callstack."""
    stack_detector = callstack_detectors.SyzyasanDetector()
    self.assertTupleEqual(
        stack_detector('Crash stack:'),
        StartOfCallStack(0, CallStackFormatType.SYZYASAN, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('A stack:'),
        StartOfCallStack(1, CallStackFormatType.SYZYASAN, LanguageType.CPP, {}))
    self.assertIsNone(stack_detector('dummy'))

  def testTsanDetector(self):
    """Tests that ``TsanDetector`` detects thread sanitizer callstack."""
    stack_detector = callstack_detectors.TsanDetector()
    self.assertTupleEqual(
        stack_detector('Read of size 1023:'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('WARNING: ThreadSanitizer'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('Previous read of size 102'),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('Location is heap block of size 3543'),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertIsNone(stack_detector('dummy'))

  def testUbsanDetector(self):
    """Tests that ``UbsanDetector`` detects ubsan callstack."""
    stack_detector = callstack_detectors.UbsanDetector()
    flag_manager = FlagManager()
    flag_manager.Register('group',
                          ParsingFlag('is_first_stack_flag', value=True))
    self.assertTupleEqual(
        stack_detector('blabla: runtime error: blabla', flag_manager),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    # After the ``is_first_stack_flag`` is set to False, the priority will be
    # 1.
    self.assertTupleEqual(
        stack_detector('blabla: runtime error: blabla', flag_manager),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertIsNone(stack_detector('dummy', flag_manager))

  def testMsanDetector(self):
    """Tests that ``MsanDetector`` detects memory sanitizer callstack."""
    stack_detector = callstack_detectors.MsanDetector()
    self.assertTupleEqual(
        stack_detector('Uninitialized value was created by'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    # After the ``is_first_stack_flag`` is set to False, the priority will be
    # 1.
    self.assertTupleEqual(
        stack_detector('Uninitialized value was stored to'),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('==123== ERROR:MemorySanitizer'),
        StartOfCallStack(2, CallStackFormatType.DEFAULT, LanguageType.CPP,
                         {'pid': 123}))
    self.assertIsNone(stack_detector('dummy'))

  def testAsanDetector(self):
    """Tests that ``AsanDetector`` detects address sanitizer callstack."""
    stack_detector = callstack_detectors.AsanDetector()
    self.assertTupleEqual(
        stack_detector('==123== ERROR:AddressSanitizer'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP,
                         {'pid': 123}))
    # After the ``is_first_stack_flag`` is set to False, the priority will be
    # 1.
    self.assertTupleEqual(
        stack_detector('READ of size 32 at backtrace:'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('freed by thread T99 here:'),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('previously allocated by thread T1 here:'),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('Thread T9 created by'),
        StartOfCallStack(1, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertIsNone(stack_detector('dummy'))

  def testChromeCrashDetector(self):
    """Tests that ``ChromeCrashDetector`` detects Fracas/Cracas callstack."""
    stack_detector = callstack_detectors.ChromeCrashStackDetector()

    self.assertTupleEqual(
        stack_detector('CRASHED [EXC @ 0x508]'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.CPP, {}))
    self.assertTupleEqual(
        stack_detector('(JAVA) CRASHED [EXC @ 0x508]'),
        StartOfCallStack(0, CallStackFormatType.DEFAULT, LanguageType.JAVA, {}))
    self.assertIsNone(stack_detector('dummy'))
