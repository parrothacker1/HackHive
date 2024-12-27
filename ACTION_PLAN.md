# The my CTFd platform

# Documentation
Swagger for API docs

## Service one 
Handle the users and teams .. CRUD of users and teams and generating JWT based on that 
* User CRUD
* Team CRUD (size set by the user(the deployer))
* Store everything in PostgreSQL
* Whenever there's a update in teamname/username notify other systems using kafka (send data to kafka)

## Service two 
Handle the creation of challenges 
* CRUD of static/dynamic challenges
* Challenge metadata like difficulty (will tend to change based on the submissions),creator and the description
* Image name will be stored in the PostgreSQL
* For static challenges we have one container made public for everyone 
* For dynamic challenges we can start/stop multiple containers based on the teams 
* Dynamic containers will be stored in Redis with the flag
* Docker API will be used to spin up the challenges
* Challenges will be divided into 2 .. one that needs a server .. and one doesn't
* Challenges will be categorized (Web,Crypto,Forencics,OSINT,Reversing,Binexp,Misc)
* Container pooling for each images that can be reused.(Where all we need to do is restart the container even after solving as it won't change any files except for the flag pooling for the rest can be done at start to avoid a overhead when challenge starts .. and static one's will not be pooled .. same for those that don't need a server (dynamic ).. for the rest .. we delete the container and make a new one)
* Container health checking
* Container status based on requests(idle,active,error .. if error add to influx .. if idle (only applicable for reusable .. for non we recreate the container) .. change the flag and reuse for another team if requesting in 3 minutes window)
* We will be using snapshots to manage containers .. 
* NOn server containers write the static file into a shared volume .. for deleting/modifying needs to aquire lock .. (redis for the locking mechanism)

## Service 3 
Handle the leaderboard and submissions
* CRUD for leaderboard
* Dynamic points based on the no.of solves and submissions  (let deployer decide the relation)
* First blood times(x as set by user)
* Whenever recieve a notification from kafka .. update the map (teamname) from postgres.


## Service 4 
Handle the admin
* There will be whitelist of IPs(if needed)
* Admin can see the submissions
* Admin can CRUD users/challenges
* Admin and moderators (admin handle CRUD of challeges and hints .. and moderators can stop accepting solutions for a challenge (admin also) but not delete)
* Moderator can ban a user but not delete
* Additional data like most solved challenges .. first blood which user etc .. 
* Whenever recieve notification , update
* when a user/team is blocked , send notification to kafka for leaderboard instance to handle

## Service 5 
Apache Kafka/RabbitMQ for notifying new challenges/hints.
We can have a websocket/ssr or something similar for this
* have admin's shedule the release of hints/challenges
* reminder for stuff like 10 min left , 1 hr left etc
* send notification if a Challenge is removed mid event

## Service 6 
Every submission will be logged into elasticsearch
And based on the submissions we will implement rate limiting on IPs

## Service 7
* InfluxDB for storing logs like response time 
* Store API calls , request paths 
* Store docker container metrics
* Retain data for at max 7 days 
* Have a script run to autoscale based on this data

## Service 8
Consul for handling normal services for service discovery for other nodes and nginx node 

## Service 9 
nginx for routing and load balancing 

## Service 10 
* file serving node for static file challenges(rev,etc ones where only static files are there in the challenge no server)
* use redis for locking
* need rate limiting
