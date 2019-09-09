# Gitlab Runner
Setup Gitlab runner and configure them to be used by GitlabCI. These runners are dynamically created in the kubernetes cluster when needed. More info [here](https://docs.gitlab.com/runner/install/kubernetes.html)

To register runners to GitlabCI you must specify the token during deployment:
```
# gitlabci/gitlab-runner/values.yaml
gitlabUrl: https://gitlab.example.com/
runnerRegistrationToken: ${RUNNER_REGISTRATION_TOKEN}
```
You can find that token on your Gitlab server, details [here](https://docs.gitlab.com/ee/ci/runners/#registering-a-shared-runner)