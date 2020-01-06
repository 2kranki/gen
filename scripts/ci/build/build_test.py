#!/usr/bin/env python3
# vi:nu:et:sts=4 ts=4 sw=4

""" Test Build.py

The module must be executed from the repository that contains the Jenkinsfile.

"""


#   This is free and unencumbered software released into the public domain.
#
#   Anyone is free to copy, modify, publish, use, compile, sell, or
#   distribute this software, either in source code form or as a compiled
#   binary, for any purpose, commercial or non-commercial, and by any
#   means.
#
#   In jurisdictions that recognize copyright laws, the author or authors
#   of this software dedicate any and all copyright interest in the
#   software to the public domain. We make this dedication for the benefit
#   of the public at large and to the detriment of our heirs and
#   successors. We intend this dedication to be an overt act of
#   relinquishment in perpetuity of all present and future rights to this
#   software under copyright law.
#
#   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
#   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
#   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
#   IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
#   OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
#   ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
#   OR OTHER DEALINGS IN THE SOFTWARE.
#
#   For more information, please refer to <http://unlicense.org/>


from unittest import TestCase
import build
import os


################################################################################
#                               Test Classes
################################################################################

class testParseArgs(TestCase):
    def setUp(self):
        pass

    def test_one(self):
        args = ["--debug"]
        build.parseArgs(args)
        self.assertTrue(build.args.debug)
        self.assertEqual(build.args.app_name, 'app01sq')


################################################################################
#                           Command-line interface
################################################################################

if '__main__' == __name__:
    import unittest
    unittest.main()

