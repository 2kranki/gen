#!/usr/bin/env python3

# vi:nu:et:sts=4 ts=4 sw=4

""" Bump the git repository version number
    Warning - This program must be executed from the repository that contains
    the 'scripts directory.
"""

import subprocess
import sys
sys.path.insert(0, './scripts')
import util                         # pylint: disable=wrong-import-position


################################################################################
#                           Object Classes and Functions
################################################################################

#---------------------------------------------------------------------
#                   Main Command Execution Class
#---------------------------------------------------------------------

class Main(util.MainBase):
    """ Main Command Execution Class
    """

    def exec_pgm(self):
        """ Execute updating the version.
        """

        if self.args.verbose > 0:
            print('*****************************************')
            print('*      Updating the Git Repo Version    *')
            print('*****************************************')
            print()

        # Read in the tag file.
        with open('scripts/tag.txt', 'r') as tag:
            ver = tag.read().strip().split('.')

        # Update the version.
        #print('.'.join(map(str, ver)))
        ver[2] = int(ver[2]) + 1
        new_ver = '.'.join(map(str, ver))
        print(new_ver)

        # Write out the new file
        if self.args.flg_exec:
            tag_out = open("scripts/tag.txt", "w")
            tag_out.write(new_ver)
            tag_out.close()

        # Now tag the git repo (git tag -a version_string -m "New Release"
        cmd = "git tag -a {0} -m \"New Release\"".format(new_ver)
        if self.args.flg_exec:
            result = subprocess.getstatusoutput(cmd)
            # result[0] == return code, result[1] == command output
            if self.args.flg_debug:
                print("\tcmd = %s" % (cmd))
                print("\trc = %s, output = %s..." % (result[0], result[1]))
            self.result_code = result[0]
        else:
            print("Would have executed:", cmd)
            self.result_code = 0
        result = subprocess.getstatusoutput("git remote")
        if int(result[0]) == 0:
            remotes = result[1]
            for remote in remotes.splitlines():
                cmd = "git push  {0} --tag".format(remote.strip())
                if self.args.flg_exec:
                    result = subprocess.getstatusoutput(cmd)
                    if self.args.flg_debug:
                        print("\tcmd = %s" % (cmd))
                        print("\trc = %s, output = %s..." % (result[0], result[1]))
                    self.result_code = result[0]
                    if result[0] != 0:
                        break
                else:
                    print("Would have executed:", cmd)
        else:
            self.result_code = result[0]

################################################################################
#                           Command-line interface
################################################################################

if __name__ == '__main__':
    Main().run()
