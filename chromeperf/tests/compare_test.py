# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import pytest

from chromeperf.pinpoint.models.compare import compare
from chromeperf.pinpoint.models.compare import thresholds


def lrange(*args):
    """Convenience helper for list(range(...)).

    Useful for generating test inputs for functions that expect concrete lists
    that can be sorted and indxed, not merely iterables.
    """
    return list(range(*args))


# Basic tests.
def test_compare_no_values_a():
    comparison = compare.compare([], [0] * 10, 10, 'functional', 1)
    assert comparison.result == compare.UNKNOWN
    assert comparison.p_value is None


def test_compare_no_values_in_b():
    comparison = compare.compare(lrange(10), [], 10, 'performance', 1)
    assert comparison.result == compare.UNKNOWN
    assert comparison.p_value is None


# Tests for compare with mode=='functional'.
def test_compare_functional_different():
    comparison = compare.compare([0] * 10, [1] * 10, 10, 'functional', 0.5)
    assert comparison.result == compare.DIFFERENT
    assert comparison.p_value <= comparison.low_threshold


def test_compare_functional_unknown():
    comparison = compare.compare([0] * 10, [0] * 9 + [1], 10, 'functional', 0.5)
    assert comparison.result == compare.UNKNOWN
    assert comparison.p_value <= comparison.high_threshold


def test_compare_functional_same():
    comparison = compare.compare([0] * 10, [0] * 10, 10, 'functional', 0.5)
    assert comparison.result == compare.SAME
    assert comparison.p_value > comparison.high_threshold


# Tests for compare with mode=='performance'.
def test_compare_performance_different():
    comparison = compare.compare(lrange(10), lrange(7, 17), 10, 'performance',
                                 1.0)
    assert comparison.result == compare.DIFFERENT
    assert comparison.p_value <= comparison.low_threshold

def test_compare_performance_unknown():
    comparison = compare.compare(lrange(10), lrange(3, 13), 10, 'performance',
                                 1.0)
    assert comparison.result == compare.UNKNOWN
    assert comparison.p_value <= comparison.high_threshold

def test_compare_performance_same():
    comparison = compare.compare(lrange(10), lrange(10), 10, 'performance', 1.0)
    assert comparison.result == compare.SAME
    assert comparison.p_value > comparison.high_threshold


def test_high_threshold_unknown_mode():
    with pytest.raises(NotImplementedError):
        thresholds.high_threshold('unknown mode', 1, 20)

def test_high_threshold_functional():
    assert thresholds.high_threshold('functional', 0.5, 20) == 0.0195

def test_high_threshold_performance():
    assert thresholds.high_threshold(
        'performance', 1.5, 20) <= thresholds.low_threshold()

def test_high_threshold_low_magnitude():
    assert thresholds.high_threshold('performance', 0.1, 20) <= 0.99

def test_high_threshold_high_magnitude():
    assert thresholds.high_threshold('performance', 10, 5) == 0.0122

def test_high_threshold_high_sample_size():
    assert thresholds.high_threshold(
        'performance', 1.5, 50) <= thresholds.low_threshold()

def test_low_threshold():
    assert thresholds.low_threshold() == 0.01
