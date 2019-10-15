#!/usr/bin/env python3

# vi:nu:et:sts=4 ts=4 sw=4

'''         Generate SQL Applications for all the Test01 Input Data

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


import      argparse
import      math
import      os
import      re
import      stat
import      subprocess
import      sys
import      time

oArgs           = None
szGenappName    = 'genapp'
szMiscDir       = './misc'
szSrcDir        = './cmd'
szTestDir       = './misc'
szTestSufixes   = ['01ma', '01ms', '01my', '01pg', '01sq']




################################################################################
#                           Object Classes and Functions
################################################################################

#---------------------------------------------------------------------
#       getAbsolutePath -- Convert a Path to an absolute path
#---------------------------------------------------------------------

def getAbsolutePath( szPath, fCreateDirs=False ):
    '''Convert Path to an absolute path.'''

    # Convert the path.
    szWork = os.path.normpath( szPath )
    szWork = os.path.expanduser( szWork )
    szWork = os.path.expandvars( szWork )
    szWork = os.path.abspath( szWork )

    if fCreateDirs:
        szDir = os.path.dirname(szWork)
        if len(szDir) > 0:
            if not os.path.exists(szDir):
                os.makedirs(szDir)

    # Return to caller.
    return szWork


def build():
    ''' Build genapp

    '''

    iRc = 0
    try:
        szCmd = 'go build -o "{0}" -v -race {1}'.format(os.path.join(oArgs.szBinDir, szGenappName),
                                                        os.path.join(szSrcDir, szGenappName, '*.go'))
        if oArgs.fDebug:
            print("\tWould have executed:", szCmd)
        else:
            if not os.path.exists(oArgs.szBinDir):
                os.makedirs(oArgs.szBinDir, 0o777)
            print("\tExecuting:", szCmd)
            os.system(szCmd)
    except OSError:
        iRc = 4

    # Return to caller.
    return iRc


def genapp(szExecFileName, szOutPath):
    ''' Generate a test application.

    :arg szExecFileName:    Exec JSON file name which is expected
                            to be in the szMiscDir.
    :arg szOutPath:         path to write the output to.
    '''

    szExecPath = os.path.join(szMiscDir, szExecFileName)

    iRc = 0
    szCmd = '"{0}/{1}" --mdldir {2} -x {3}'.format(oArgs.szBinDir, szGenappName, oArgs.szModelDir, szExecPath)
    try:
        if oArgs.fDebug:
            print("\tWould have executed:", szCmd)
        else:
            print("\tExecuting:", szCmd)
            os.system(szCmd)
    except OSError:
        iRc = 4

    if iRc == 0:
        iRc = genJenkins(szOutPath)

    # Return to caller.
    return iRc


def         parseArgs(listArgV=None):
    '''
    '''
    global      oArgs

    # Parse the command line.
    szUsage = "usage: %prog [options] sourceDirectoryPath [destinationDirectoryPath]"
    oCmdPrs = argparse.ArgumentParser( )
    oCmdPrs.add_argument('-b', '--build', action='store_false', dest='fBuild',
                         default=True, help='Do not build genapp before using it'
                         )
    oCmdPrs.add_argument('-d', '--debug', action='store_true', dest='fDebug',
                         default=False, help='Set debug mode'
                         )
    oCmdPrs.add_argument('-f', '--force', action='store_true', dest='fForce',
                         default=False, help='Set force mode'
                         )
    oCmdPrs.add_argument( '-v', '--verbose', action='count', default=1,
                        dest='iVerbose', help='increase output verbosity'
                        )
    oCmdPrs.add_argument('--appdir', action='store', dest='szAppDir',
                         default='/tmp', help='Set Application Base Directory'
                         )
    oCmdPrs.add_argument('--appname', action='store', dest='szAppName',
                         default='app01', help='Set Application Base Name'
                         )
    oCmdPrs.add_argument('--bindir', action='store', dest='szBinDir',
                         default='/tmp/bin', help='Set Binary Directory'
                         )
    oCmdPrs.add_argument('--mdldir', action='store', dest='szModelDir',
                         default='./models', help='Set genapp Model Directory'
                         )
    oCmdPrs.add_argument('args', nargs=argparse.REMAINDER, default=[])
    oArgs = oCmdPrs.parse_args( listArgV )
    if oArgs.iVerbose:
        print('*****************************************')
        print('*      Generating the Application       *')
        print('*****************************************')
        print()
    oArgs.szAppPath = os.path.join(oArgs.szAppDir, oArgs.szAppName)
    if oArgs.fDebug:
        print("In DEBUG Mode...")
        print('Args:', oArgs)




################################################################################
#                           Main Program Processing
################################################################################

def         mainCLI(listArgV=None):
    '''Command-line interface.'''
    execFlags = (stat.S_IRUSR | stat.S_IWUSR | stat.S_IXUSR | stat.S_IRGRP |
                 stat.S_IRWXO | stat.S_IROTH | stat.S_IXOTH)
    global      szAppPath
    global      oArgs
    global      oLogOut

    # Do initialization.
    iRc = 20

    parseArgs(listArgV)

    if len(oArgs.args) < 1:
        szSrc = os.getcwd( )
    else:
        szSrc = oArgs.args[0]
    if len(oArgs.args) > 1:
        print("ERROR - too many command arguments!")
        oCmdPrs.print_help( )
        return 4
    if oArgs.fDebug:
        print('szSrc:',szSrc)

    # Set up base objects, files and directories.
    if not os.path.exists(oArgs.szAppPath):
        print("\tCreating Directory:",oArgs.szAppPath)
        os.makedirs(oArgs.szAppPath)

    # Perform the specified actions.
    iRc = 0
    try:

        # Build genapp if needed.
        if oArgs.fBuild:
            print("\tBuilding genapp...")
            build()

        # Generate the application subdirectories.
        for v in szTestSufixes:
            print("\tCreating app for app{0}...".format(v))
            szPath = os.path.join(oArgs.szAppPath, "app{0}".format(v))
            iRc = genapp("test{0}.exec.json.txt".format(v), szPath)
            if iRc != 0:
                break

    finally:
        pass
    print()

    return iRc




################################################################################
#                           Command-line interface
################################################################################

if '__main__' == __name__:
    startTime = time.time( )
    iRc = mainCLI( sys.argv[1:] )
    if oArgs.iVerbose or oArgs.fDebug:
        if 0 == iRc:
            print("...Successful completion.")
        else:
            print("...Completion Failure of %d" % (iRc))
    endTime = time.time( )
    if oArgs.iVerbose or oArgs.fDebug:
        print("Start Time: %s" % (time.ctime(startTime)))
        print("End   Time: %s" % (time.ctime(endTime)))
        diffTime = endTime - startTime      # float Time in seconds
        iSecs = int(diffTime % 60.0)
        iMins = int((diffTime / 60.0) % 60.0)
        iHrs = int(diffTime / 3600.0)
        print("run   Time: %d:%02d:%02d" % (iHrs, iMins, iSecs))
    sys.exit( iRc or 0 )


