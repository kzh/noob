### Noob
A scalable distributed microservice architecture platform running on Kubernetes capable of judging and orchestrating the execution of user submitted code.

<img src="https://raw.githubusercontent.com/kzh/noob/8e83d650de4d94855582c615944ef166819081f9/images/Noob%20Architecture.svg?sanitize=true">

#### Roadmap
- Set up Kubernetes Cluster.
    * Set up rook helm chart for persistent volumes. (✅ 8/5/18)
        * ^ decided to do this w/o Helm.
    * Set up redis helm chart for sessions. (✅ 8/5/18)
    * Set up mongodb helm chart. (✅ 8/5/18)
    * Se up rabbitmq helm chart for code execution queue.
- Set up auth microservice.
    * Add Dockerfile. (✅ 8/10/18)
    * Integrate into Kubernetes/Helm:
        * Create deployment. (✅ 8/11/18)
        * Create service. (✅ 8/11/18)
        * Set up /auth/ with nginx ingress. (✅ 8/19/18)
    * Connect to redis. (✅ 8/12/18)
    * Connect to mongodb. (✅ 8/13/18)
    * Endpoints:
        * Login - authenticate username + password with mongodb,              create session in redis.               (✅ 8/17/18)
        * Logout - destroy session in redis. (✅ 8/17/18)
        * Register - store username + hashed password in mongodb. (✅ 8/17/18)
- Set up frontend microservice.
    * Add Dockerfile. (✅ 8/28/18)
    * Integrate into Kubernetes/Helm:
        * Create deployment. (✅ 8/28/18)
        * Create service. (✅ 8/28/18)
        * Set up / with nginx ingress. (✅ 8/28/18)
    * Create mock authentication page for signing up and logging in. (✅ 8/28/18)
- Set up admin microservice.
    * Add Dockerfile. (✅ 9/5/18)
    * Integrate into Kubernetes/Helm:
        * Create deployment. (✅ 9/5/18)
        * Create service. (✅ 9/5/18)
        * Set up /admin/ nginx ingress. (✅ 9/5/18)
    * Endpoints:
        * Create Problem - store problem in mongodb (✅ 9/5/18)
        * Update Problem - update problem in mongodb (✅ 9/5/18)
        * Delete Problem - delete problem in mongodb (✅ 9/5/18)
    * Frontend UI
- next: tbd…

### Development
Some commands to know :P .  
**Accessing MongoDB**:
```
$ kubectl get secret noob-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode
$ kubectl run -i -t --rm debug --image=ubuntu --restart=Never
$ apt-get update && apt-get install mongodb
$ mongo --host noob-mongodb --port 27017 -u root -p <PASSWORD> admin
```
   * Changing user’s role:
```
db.users.findAndModify({
   query: {username: <USERNAME>},
   update: {$set: {role: <ROLE>}},
   new: true,
})
```

**Updating Microservices**:
* Entire System:
```
$ docker-compose build && docker-compose push
$ helm delete --purge noob
$ helm install --namespace noob --name noob ./chart/
```
* Single Microservice:
```
$ docker-compose build <microservice> && docker-compose push <microservice>
$ kubectl get pods
$ kubectl delete pod <microservice>
```

#### Personal Notes
- Helm update will mess up the redis k8s secret since the secret does not update while the redis password will. The solution to this is just do a hard delete and install when updating the entire chart. (Temp fixed)

#### Sidetrack
- Set up continuous integration and deployment.
    * Possibly with Google Cloud Build or Concourse?
- Make sure to have high quality documentation!
- Also test, test, test!
- Deploy to Google Cloud Platform or Amazon Web Services. Currently running my own Kubernetes cluster on OVH vps.
