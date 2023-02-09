
# AnaLog

Note: This Project is a work in progress. Help is appreciated! Look at the contributing section to look at some features that could be worked on!

Tool to mine rules to reduce human effort in parsing logs. No more random regex searches! 

Uses frequency based phrase mining and a static "database" of common regular expressions to create clues that help developers dig into a mountain of logs with the least effort possible. 







## Badges

[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)


## Authors

- [@adiraokhoury](https://www.github.com/adiraokhoury)


## Contributing

Contributions are always welcome!
Feature Requests are tracked in [this](https://github.com/users/adiraokhoury/projects/5) github project board. 



## Blog Post


The internals of why and how I built this is detailed in [WORK IN PROGRESS] on my [website](https://adiraokhoury.github.io). 
## Lessons Learned


This project was a beginners attempt to use Go for text processing and file handling. 

Using golang is an absolute pain when it comes to declaring nested hashmaps, and understandably so. The tradeoff is in how easy it is to define something as a composition of structures, rather than a complicated dictionary object using only base types. 

The reflection package made it a lot easier to access the properties of a specific instance of some structure. 

More to be added in due time. 



## Who is this for? 

Developers who have to deal with a large volume of logs, and have no idea where to start. 
