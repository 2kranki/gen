                        Gen Projects

The idea for the Gen project came about when I became aware of 
how powerful the templating system in Go is.  So, as I was taking
the udemy.com course, "Web Development with Google's Golang", I
decided to write a program generator for simple SQL Table websites.

I want the generated code to stand on its own and provide a starting
point for further development.  I do not see the generation process
as the end all. So, the generation process must output good enough
quality code that I would not mind working on that code.  

The system is based on defining hjson files for the initial application
definition and having them control the generation process. (I have
previously written an hjson parser in C and use it in my C applications
and libraries.)

It is still a work in progress and needs more work which I am in
the progress of adding.  However, some of the fundamentals are done
and I will build and refine them as I go. This project is still around
version 0.07.  My current road map is to:

1. Finish handlers_test.go.
2. Get the newly generated application fully operational.
3. Finish main_test.go
4. Declare version 0.1

