### Noob [August 2018]

The goal is to create a leetcode clone on Kubernetes. Noob is mainly for me to learn the popular cloud tools being Kubernetes, Docker, etc as well as to explore microservice architecture and further develop my system design experience. The hope is for Noob to be scalable so in case of increasing traffic, Noob can easily adapt and handle the load.

#### Roadmap
- Set up Kubernetes Cluster.
    * Set up rook helm chart for persistent volumes. (✅ 8/5/18)
        * ^ decided to do this w/o Helm.
    * Set up redis helm chart for sessions. (✅ 8/5/18)
    * Set up mongodb helm chart. (✅ 8/5/18)
- Set up auth microservice.
    * Add Dockerfile. (✅ 8/10/18)
    * Integrate into Kubernetes/Helm:
        * Create deployment. (✅ 8/11/18)
        * Create service. (✅ 8/11/18)
        * Set up /api/auth with nginx ingress.
    * Connect to redis. (✅ 8/12/18)
    * Connect to mongodb. (✅ 8/13/18)
    * Endpoints:
        * Login - authenticate username + password with mongodb,              create session in redis.               (✅ 8/17/18)
        * Logout - destroy session in redis. (✅ 8/17/18)
        * Register - store username + hashed password in mongodb. (✅ 8/17/18)
- Set up frontend microservice.
    * Integrate into Kubernetes/Helm:
        * Create deployment.
        * Create service.
        * Set up nginx ingress.
    * Set up nginx http server to serve static files.
    * Create mock authentication page for signing up and logging in.
- next: tbd…

#### Personal Notes
- wip

#### Sidetrack
- Set up continuous integration and deployment.
    * Possibly with Google Cloud Build or Concourse?
- Make sure to have high quality documentation!
- Also test, test, test!
- Deploy to Google Cloud Platform or Amazon Web Services. Currently running my own Kubernetes cluster on OVH vps.
