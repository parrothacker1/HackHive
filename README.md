# HackHive

```
The Project is still under development.
```

Welcome to HackHive. A CTF framework (more precisely just the backend) under development which aims to make a ready-to-go CTF platform for your event. This platform will support both static and dynamic flags and uses Docker for spinning up the challenge.

## Contents
* Getting Started
* Roadmap
* Contributions
* Future Plans

## Getting Started

To get started with working on the project, first you will have to clone the repository.To clone the repository, run 

```
git clone https://github.com/parrothacker1/HackHive
```

To start coding and working, I would suggest installing rust, golang and python (for now this is subjected to change depending on the code and the web framework that suits the service best).Also I would recommend using Docker for building the image as all of ours services will be dockerized. (PS: If you have docker, you might be able to avoid installing python and the rest as you can directly compile on it.)

To install docker, follow their [docs](https://docs.docker.com/engine/install/).<br>
To install rust, follow their [docs](https://www.rust-lang.org/tools/install).<br>
To install golang, follow their [docs](https://go.dev/doc/install).<br>
To install python, follow their [docs](https://www.python.org/downloads/) (although I don't think you need to refer that).

After cloning the repository , checkout the Roadmap to understand how we are designing the backend.


## Roadmap

### Service 1 (Users)
This service mainly focuses on CRUD of users and teams and generation of JWT for other services to authorize and authenticate users. We are using [net/http](https://pkg.go.dev/net/http) for handling requests.For those who are wondering why not something better like fasthttp, well this service's main job is creating users and handling it's related work like updating,deleting etc and I find this framework good for this. And as valyala said, [fasthttp](https://github.com/valyala/fasthttp) is meant for an edgecase where you want to handle thousands of requests and need a consistent low millisecond response time.

### Service 2 (Containers)
This service mainly focuses on spinning up containers for each challenges.For this service we are planning to use [Actix](https://actix.rs) as for this process I believe we need the raw power of rust to handle things.<br>
(If you have a better plan please do tell me)

### Service 3 (Admin)
I still haven't figured out what framework to use here, but the admin's main work is to monitor the users and challenges and get the demographic of how frequent a challenge is attempted and etc stuff like that.

### Service 4 (LeaderBoard)
This service's main work is to get all the submissions and update the leaderboard. Here we are using [fasthttp](https://github.com/valyala/fasthttp) for handling all the requests and I plan to store the updated data in Redis.

### Service 5 (FileShare)
WE plan to use [FastAPI](https://fastapi.tiangolo.com/) for sharing the files. So the files here are the static files related to the challenge. Like for example, if a challenge has a binary file which needs to be reversed, we can share that file through this service. Instead of having the container host these files, we can have the container generate the necessary files and store it in disk and use this service to get the file. This way we can free the container after it's use and we have more space to run the container

**There's much more services** like Apache Kafka(for notifying the users about the new challenges, time left,etc) and ElasticSearch (for logging purposes), Redis (for caching), Nginx and Consul for scaling.Although these services will be setup in the later part of development. I will write a detailed description of what these services will do.


## Contributions

To contribute to this project, check out our CONTRIBUTING.md file for more details. This will help you to understand what we are prioritizing right now and what we would really appreciate. 😃

## Future Plans
So right now we are making a huge framework with dynamic and static flags. The future plans of this project include :
* Adding support for (AWD) Attack and Defense
* Adding support for (KOTH) King of the Hill 
* Creating a small scale framework for just hosting static challenges.

The framework is still incomplete and I believe there's more space for optimization and more features I can incoperate in this.

## Support this Project
This project is developed with passion and dedication to make a really good framework with support for all type of CTF events. If you find it useful and want to support its ongoing development, please consider making a donation. Your contributions help keep the project alive and growing. 🌱
If you find this project useful,you can show your support in a few ways:
* 🌟 Star this repository: Starring helps others discover the project and boosts its visibility.
* 📢 Share it with your friends: Spread the word to those who might benefit from or contribute to this project.
* ☕ Make a donation: If you'd like to support ongoing development and future enhancements, consider donating:
    - [Buy Me a Coffee](https://buymeacoffee.com/parrothacker1) 

Every contribution, whether it's a donation, a star, or simply sharing the project, makes a huge difference. Thank you for your support! 🙏
