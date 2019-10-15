#!/usr/bin/env python3

# vi:nu:et:sts=4 ts=4 sw=4

'''       Test Generate SQL Applications for all the Test01 Input Data

        Test01 Input Data has test data for each SQL Server type supported
        by genapp so that it can be properly tested. This program scans
        ./misc/ for all the test01 application definitions and generates
        them.  Included in the generation process is the generation of
        Jenkins support for building, testing and deploying to Docker Hub
        each of the applications.

        To simplify things, this script must be self-contained using only
        the standard Python Library.


TODO:
    -   Add some more tests

'''

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


from io import StringIO
from unittest import TestCase
import genSql01
import os


################################################################################
#                               Test Classes
################################################################################

class testParseArgs(TestCase):
    def setUp(self):
        pass

    def test_one(self):
        args = ["--debug"]
        genSql01.parseArgs(args)
        self.assertTrue(genSql01.oArgs.fDebug)


class testAbsolutePath(TestCase):
    def setUp(self):
        genSql01.parseArgs(["--debug"])

    def test_one(self):
        txt = "${HOME}/a.txt"
        a = genSql01.getAbsolutePath(txt)
        b = os.path.expandvars(txt)
        self.assertEqual(a,b)
        a = genSql01.getAbsolutePath('~/a.txt')
        self.assertEqual(a,b)


class testBuild(TestCase):
    def setUp(self):
        genSql01.parseArgs(["--debug"])

    def test_one(self):
        iRc = genSql01.build()
        self.assertEqual(iRc, 0)


class testJenkins(TestCase):
    def setUp(self):
        genSql01.parseArgs(["--debug"])

    def test_one(self):
        iRc = genSql01.genJenkins("/tmp/app01/app01ma")
        self.assertEqual(iRc, 0)


################################################################################
#                           Command-line interface
################################################################################

if '__main__' == __name__:
    import unittest
    unittest.main()

