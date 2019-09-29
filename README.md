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
* Supports MariaDB, MS SQL Server, MySQL, PostGres and SQLite.
* In testing app01, Customer table works in all sql servers, Vendor table works
    in all but SQLite for now.


**Project Futures**:


- [x] Finish Table Handler Testing
- [x] Finish Database I/O Testing for SQLite
- [x] Finish Table I/O Testing for SQLite
- [x] Finish Database Handler Testing
- [x] Finish MS-SQL Support
- [x] Finish MySQL/MariaDB Support
- [x] Finish Postgres Support
- [x] Change structure of generated application to conform to Go's package requirements. 
- [x] Changed genapp to conform to Go's new module format.
- [x] Clean up form templating
- [x] Add more testing to the generated code
- [x] Add **csv** save/restore per table
- [x] Convert golang compile and testing to use Docker so that it is the same for linux, 
        macOS and windows if possible. Right now, the scripts are for bash and only
        support linux and macOS
- [x] Add docker-compose yml script to run the MariaDB, MS Server, MySQL and PostGres generated programs.
- [x] Declare version 0.2
- [x] Fix sqlite generated code to run in Docker
- [x] Changed all displays to contain a possible completion message
- [x] Change structure of application to conform to Go's module/package requirements/standards.
- [ ] Change handlers that issue text message to issue html page with message
- [ ] Change haneler test routines to parse generated html and check for data
- [ ] Add HTTPS support 
- [ ] Add User Authentication and x- header authentication support
- [ ] Add support for JSON data to/from HTTP client
- [ ] Add JSON analysis phase that looks for errors in the definitions ahead of
        code generation such as SQLite rowid analysis.
- [ ] Declare version 0.3
- [ ] Add GORM support 
- [ ] Add SQL Generation for C#/.Net applications 


**Usage**:


To use this as I do, try the following:
1. Install GO
2. `git clone github.com/2kranki/genapp.git`
3. `cd genapp`
4. `./b.sh`     <- builds /tmp/bin/genapp
5. `./gen01.sh` will generate /tmp/app01 for all the sql servers.
6. `cd /tmp/app01/app01xx` where xx stands for ma (MariaDB), ms (MS SQL), my (MySQL), pg (PostGres), or sq (SQLite)
7. `./b.sh`     <- builds app01 as app
8. `/tmp/bin/app01xx` <- (See Step 6 for xx) SQLite Application on localhost:8080
9. `localhost:8080` in your browser
10. Select "Load test data" for Customer Table
11. Back to main menu
12. Select "Maintain Table" for Customer Table
13. Play and have fun.
14. Please send me any comments or problems.

Look in the "dbs" directory for specific notes on what I did to get each database driver running.  Each was a little different on my system and it might be that way for you.  Remember that you can over-ride the connection parameters from the command line.  To see the arguments, just run "/tmp/bin/app --help" and it will display them.  Actually, I no longer use the dbs shell scripts much. They should still work. I am just migrating to Python scripts, Docker and Jenkins.

For now, I generate Docker images using ./jenkins/build/build.sh. Then I cd to the appropriate subdirectory and run "docker-compose up &" plus "docker-compose down". By executing this way, I can see all the tracing that is built-in (ie container log) on the terminal and execute the program from my browser to test it.


I am running this on MacOSX for now with golang v1.12.7.  I will adapt the above for Windows at some point.
My editor is 'Goland' from Intellij which I like.

I have added repository, app01. It contains the latest output of *genapp* for the Test01 example.  See the *misc* directory in this repo for the controlling JSON files and [app01](https://github.com/2kranki/app01) repostiory for the generated code. The README in [app01](https://github.com/2kranki/app01) explains how to run the generated applications as well. *gen01.sh* generates app01 for all sql servers supported.

If you are new to Golang, I strongly recommend that you buy "The GO Programming Language" by
Alan Donovan and Brian Kernighan.  It is an excellent book and reference.  I also bought "GO in Action" by William Kennedy but I can not recommend it given it's cost compared to "The GO Programming Language" book.  "The GO Programming Language" book is far superior from my point of view and is a great reference on the language and its libraries.

Last, I apologize for the ever changing structure and API. I started this project to learn GO. On the way, I am continuing to learn new techniques and layouts. I have not strived to maintain any backwards consistency. You will find this project to change more. I am studying further techniques in Docker, RESTful APIs and GO application structure. I am constanly learning new things everyday. It is wonderful!
