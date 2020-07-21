# Knative - Kubernetes Job Trigger

[![Go Report Card](https://goreportcard.com/badge/github.com/itsmurugappan/job-trigger)](https://goreportcard.com/report/github.com/itsmurugappan/job-trigger)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-59%25-brightgreen.svg?longCache=true&style=flat)</a>

KNative Service for creating kubernetes jobs

### PreReq

Please create the below rolebinding to give access to your ksvc to create a job

```
kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: default-rb
  namespace: <ns-name>
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
- kind: ServiceAccount
  name: default
  namespace: <ns-name>
EOF
```

### Creating and Running the service

This knative service looks for the job spec in a environment variable named `spec`.
Below is an example of a service specification

```
kubectl apply -f - <<EOF
apiVersion: serving.knative.dev/v1alpha1
kind: Service
metadata:
  name: sample-job-trigger-svc
spec:
  template:
    spec:
      containers:
      - env:
        - name: spec
          value: "{\"Image\": \"murugappans/goswaggertest\",\"Name\": \"sample-job\"}"
        image: ko://github.com/itsmurugappan/job-trigger/cmd/job-trigger
EOF     
```

### Configuring your Job

#### CMD

```
  "Name": "spec",
  "Value": "{\"Image\": \"murugappans/goswaggertest\",\"Name\": \"simplego\",\"User\": 1002, \"Cmd\": [\"python\", \"job.py\"]}"
```

#### Environment Variables

```
  "Name": "spec",
  "Value": "{\"Image\": \"murugappans/goswaggertest\",\"Name\": \"simplego\",\"User\": 1002, \"EnvVariables\": [\"Name\": \"env\",\"Value\": \"prod\"}]}"
```

#### CM and Secret volumes 

```
  "Name": "spec",
  "Value": "{\"Image\": \"murugappans/goswaggertest\",\"Name\": \"simplego\",\"User\": 1002, \"Secrets\": [{\"Name\": \"secretname\",\"MountPath\": \"pathtomount\"}],\"ConfigMaps\": [{\"Name\": \"cmname\",\"MountPath\": \"pathtomount\"}]}"
```

####  CM and Secret as env variables

```
  "Name": "spec",
  "Value": "{\"Image\": \"murugappans/goswaggertest\",\"Name\": \"simplego\",\"User\": 1002, \"EnvFromSecretorCM\": [{\"Name\": \"secretname\",\"Type\": \"Secret\"}]}"
```


### Dynamic input

At run time you might want to inject some variables. Those variables can be passed as url query parameters

For example if the function url is as below. Key and Value will be mounted as environment variable in the job

```
http://sample-job-trigger.test.com?key=value
```

### Labeling your Job

In some cases you might want to label your job. It can be provided in the query string as below.

```
http://sample-job-trigger.test.com?labels=label1=value,label2=value
```

### CleanUp

To clean up old jobs, specify the `history` query parameter. This param specifies the number of job logs to be preseverd. 

For example if history is 0, all logs will be will be deleted, for 2 , 2 latest will be preserved.

Default value is 3

```
http://sample-job-trigger.test.com?history=1
```
