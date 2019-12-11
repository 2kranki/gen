#!/usr/bin/env python3

# vi:nu:et:sts=4 ts=4 sw=4

"""         Generate SQL Applications for all the Test01 Input Data

        Test01 Input Data has test data for each SQL Server type supported
        by genapp so that it can be properly tested. This program scans
        ./misc/ for all the test01 application definitions and generates
        them.  Included in the generation process is the generation of
        Jenkins support for building, testing and deploying to Docker Hub
        each of the applications.

        To simplify things, this script must be self-contained using only
        the standard Python Library.


TODO:
    -   Finish jenkins/build generation
    -   Finish jenkins/deploy generation
    -   Finish jenkins/push generation
    -   Finish jenkins/test generation
    -   Finish jenkinsfile generation
    -   Finish application generation

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


import os
import subprocess
import sys
sys.path.insert(0, './scripts')
import util                         # pylint: disable=wrong-import-position


################################################################################
#                           Object Classes and Functions
################################################################################

class Main(util.MainBase):
    """ Main Command Execution Class
    """

    def __init__(self):
        super().__init__()
        self.test_suffixes = ['01ma', '01ms', '01my', '01pg', '01sq']
        self.genapp_name = 'genapp'

    def arg_parse_add(self):
        """ Add additional arguments.
        """
        self.arg_prs.add_argument('-b', '--build', action='store_false', dest='flg_build',
                                  default=True, help='Do not build genapp before using it'
                                 )
        self.arg_prs.add_argument('--appdir', action='store', dest='app_dir',
                                  default='/tmp', help='Set Application Base Directory'
                                 )
        self.arg_prs.add_argument('--appname', action='store', dest='app_name',
                                  default='app01', help='Set Application Base Name'
                                 )
        self.arg_prs.add_argument('--bindir', action='store', dest='bin_dir',
                                  default='/tmp/bin', help='Set Binary Directory'
                                 )
        self.arg_prs.add_argument('--srcdir', action='store', dest='src_dir',
                                  default='./cmd', help='Set genapp Source Directory'
                                 )
        self.arg_prs.add_argument('--mdldir', action='store', dest='mdl_dir',
                                  default='./models', help='Set genapp Model Directory'
                                 )
        self.arg_prs.add_argument('--mscdir', action='store', dest='msc_dir',
                                  default='./misc', help='Set genapp Misc Directory'
                                 )
        self.arg_prs.add_argument('--tstdir', action='store', dest='tst_dir',
                                  default='./misc', help='Set genapp Test Directory'
                                 )

    def arg_parse_exec(self):
        """ Execute the argument parsing.
            Warning - Main should override this method if additional cli
            arguments are needed or argparse needs some form of modification
            before execution.
        """
        self.arg_parse_setup()
        self.arg_parse_add()
        self.arg_parse_parse()
        self.args.app_path = os.path.join(self.args.bin_dir, self.genapp_name)

    def build(self):
        """ Build the Golang program, genapp.
        """
        try:
            src_path = os.path.join(self.args.src_dir, self.genapp_name, '*.go')
            if self.args.flg_debug:
                print("\tapp_path:", self.args.app_path)
                print("\tsrc_path:", src_path)
            cmd = 'go build -o "{0}" -v -race {1}'.format(self.args.app_path, src_path)
            if not self.args.flg_exec:
                print("\tWould have executed:", cmd)
                self.result_code = 0
            else:
                if not os.path.exists(self.args.bin_dir):
                    os.makedirs(self.args.bin_dir, 0o777)
                print("\tExecuting:", cmd)
                self.result_code = subprocess.call(cmd, shell=True)
        except:                                         # pylint: disable=bare-except
            self.result_code = 8

    def exec_pgm(self):                                 # pylint: disable=no-self-use
        """ Program Execution
            Warning - Main should override this method and make certain that
            it returns an exit code in self.result_code.
        """
        if len(self.args.args) > 0:
            print("ERROR - too many command arguments!")
            self.arg_prs.print_help()
            self.result_code = 0
            return
        if self.args.flg_debug:
            print('\tsrc_dir:', self.args.src_dir)

        # Set up base objects, files and directories.
        if not os.path.exists(self.args.app_path):
            print("\tCreating Directory:", self.args.app_path)
            if self.args.flg_exec:
                os.makedirs(self.args.app_path)
            else:
                print("\tWould have executed: mkdir -p", self.args.app_path)

        # Perform the specified actions.
        try:
            # Build genapp if needed.
            if self.args.flg_build:
                print("\tBuilding genapp...")
                self.build()
            # Generate the application subdirectories.
            for suffix in self.test_suffixes:
                print("\tCreating app for app{0}...".format(suffix))
                self.genapp("test{0}.exec.json.txt".format(suffix))
                if self.result_code != 0:
                    break
        finally:
            pass
        print()

    def genapp(self, file_name):
        """ Generate a test application.
            :arg szExecFileName:    Exec JSON file name which is expected
                                    to be in the szMiscDir.
            :arg szOutPath:         path to write the output to.
        """
        exec_path = os.path.join(self.args.msc_dir, file_name)
        app_path = os.path.join(self.args.bin_dir, self.genapp_name)
        cmd = '"{0}" --mdldir {1} -x {2}'.format(app_path, self.args.mdl_dir, exec_path)
        try:
            self.result_code = 0
            if self.args.flg_exec:
                print("\tExecuting:", cmd)
                os.system(cmd)
            else:
                print("\tWould have executed:", cmd)
        except:                                         # pylint: disable=bare-except
            self.result_code = 4


################################################################################
#                           Command-line interface
################################################################################

if  __name__ == '__main__':
    Main().run()
