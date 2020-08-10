# Dokku Deployment

- Add .gitlab-ci.yml:
```yaml
stages:
  - deploy
  
variables:
  APP_NAME: betting-game
  
deploy:
  stage: deploy
  image: ilyasemenov/gitlab-ci-git-push
  environment:
    name: production
  only:
    - master
  script:
    - git-push ssh://dokku@$VM_IP:22/$APP_NAME
```

- Create project on Gitlab
- Add remote repository to the project
- Set env vars VM_IP (IP of VM that Dokku is installed) and SSH_PRIVATE_KEY on Gitlab
- Push project to master (`git push origin master`)

```