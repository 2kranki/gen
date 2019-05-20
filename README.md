                        *genapp* - Go SQL Application Generator - v0.1.0

The idea for the *genapp* project came about when I became aware of how powerful the templating system in Go is.  So, as I was taking the udemy.com course, "Web Development with Google's Golang" by Todd McLeod, I decided to write a program generator for simple SQL Table websites.

I want the generated code to stand on its own and provide a starting
point for further development.  I do not see the generation process
as the end all. So, the generation process must output good enough
quality code that I would not mind working on that code.  

The system is based on defining json/hjson files for the initial application
definition and having them control the generation process.


**Project Status**:


* Version is 0.1
* Supports only SQLite
* Still a work in progress


**Project Futures**:


- [ ] Finish MS-SQL Support
- [ ] Finish MySQL/MariaDB Support
- [ ] Finish Postgres Support
- [ ] Clean up form templating
- [ ] Add more testing to xx_test.go files in both genapp and the generated code
- [ ] Add **csv** save/restore per table
- [ ] Declare version 0.2


**Usage**:


To use this as I do, try the following:
1. Install GO
2. `git clone github.com/2kranki/genapp.git`
3. `cd genapp`
4. `./b.sh`     <- builds /tmp/genapp
5. `/tmp/genapp -x misc/test01.exec.json.txt` <- Creates "/tmp/app01"
6. `cd /tmp/app01/app`
7. `./b.sh`     <- builds app01 as app
8. `./app`      <- SQLite Application on localhost:8080
9. `localhost:8080` in your browser
10. Select "Load test data" for Customer Table
11. Back to main menu
12. Select "Maintain Table" for Customer Table
13. Play and have fun.
14. Please send me any comments or problems.



I am running this on MacOSX for now with golang v1.12.1.  I will adapt the above for Windows at some point.
My editor is Goland from Intellij.  It has worked well for me.

I have added repository, app01. It contains the latest output of *genapp* for the Test01
example.  See the *misc* directory in this repo for the controlling JSON files and 
[app01](https://github.com/2kranki/app01) repostiory for the generated  code.

