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
- Push project to master
- Set env vars VM_IP and SSH_PRIVATE_KEY on Gitlab


```