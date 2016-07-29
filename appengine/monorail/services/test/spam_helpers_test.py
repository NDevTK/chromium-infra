# coding=utf-8

import unittest

from services import spam_helpers

class SpamHelpersTest(unittest.TestCase):

  def testExtractUrls(self):
    urls = spam_helpers._ExtractUrls('')
    self.assertEquals(0, len(urls))
    urls = spam_helpers._ExtractUrls('check this out: http://google.com')
    self.assertEquals(1, len(urls))
    self.assertEquals(['http://google.com'], urls)

  def testEmailIsSketchy(self):
    self.assertFalse(spam_helpers._EmailIsSketchy('', ()))
    self.assertFalse(
        spam_helpers._EmailIsSketchy('jan1990@example.com', ('@example.com')))
    self.assertTrue(
        spam_helpers._EmailIsSketchy('jan1990@foo.com', ('@example.com')))

  def testHashFeatures(self):
    hashes = spam_helpers._HashFeatures(('', ''), 5)
    self.assertEqual(5, len(hashes))
    self.assertEquals([1.0, 0, 0, 0, 0], hashes)

    hashes = spam_helpers._HashFeatures(('abc', 'abc def'), 5)
    self.assertEqual(5, len(hashes))
    self.assertEquals([0, 0, 0.6666666666666666, 0, 0.3333333333333333],
        hashes)

  def testGenerateFeatures(self):
    features = spam_helpers.GenerateFeatures(
        'abc', 'abc def http://www.google.com http://www.google.com',
        'unused@chromium.org', 5, ('@chromium.org'))
    self.assertEquals(12, len(features))
    self.assertEquals(['False', '2', '1', '3', '11', '51', '39',
        '0.363636', '0.000000', '0.181818', '0.000000', '0.454545'], features)

    features = spam_helpers.GenerateFeatures(
        'abc', 'abc def', 'jan1990@bar.com', 5, ('@example.com'))
    self.assertEquals(12, len(features))
    self.assertEquals(['True', '0', '0', '3', '11', '7', '15',
        '0.000000', '0.000000', '0.666667', '0.000000', '0.333333'], features)

    # BMP Unicode
    features = spam_helpers.GenerateFeatures(
        u'abc’', u'abc ’ def', 'jan1990@bar.com', 5, ('@example.com'))
    self.assertEquals(12, len(features))
    self.assertEquals(['True', '0', '0', '6', '14', '11', '19',
        '0.000000', '0.000000', '0.250000', '0.250000', '0.500000'], features)

    # Non-BMP Unicode
    features = spam_helpers.GenerateFeatures(
        u'abc國', u'abc 國 def', 'jan1990@bar.com', 5, ('@example.com'))
    self.assertEquals(12, len(features))
    self.assertEquals(['True', '0', '0', '6', '14', '11', '19',
        '0.000000', '0.000000', '0.250000', '0.250000', '0.500000'], features)

    # A non-unicode bytestring containing unicode characters
    features = spam_helpers.GenerateFeatures(
        'abc…', 'abc … def', 'jan1990@bar.com', 5, ('@example.com'))
    self.assertEquals(12, len(features))
    self.assertEquals(['True', '0', '0', '6', '14', '11', '19', 
        '0.250000', '0.000000', '0.250000', '0.250000', '0.250000'], features)

    # Empty input
    features = spam_helpers.GenerateFeatures(
        '', '', 'jan1990@bar.com', 5, ('@example.com'))
    self.assertEquals(12, len(features))
    self.assertEquals(['True', '0', '0', '0', '8', '0', '8',
        '1.000000', '0.000000', '0.000000', '0.000000', '0.000000'], features)
