# Board Gitlab Helper

When starting Board with Gitlab service integration, this helper could assit to accomplish Gitlab pre-settings.

*By default this helper should be as Docker image that would be executed when starting Board by using ```install.sh```*

### To build image with code updates
1 Change to the directory:
   
   ```sh
   cd board/tools/gitlab-helper
   ```

2 Execute shell command.
   
   ```sh
   docker build -f container/Dockerfile -t gitlab-helper:<tag-version> .
   ```
### To run helper image manually

1 You should have ```board.cfg``` file changed expectly.***Please DO NOT edit ```GITLAB_ACCESS_TOKEN``` by yourself in this config file.***

2 **This image can be only effected when you have not started a Gitlab service yet.**

3 Execute shell command.
   ```
   docker run --rm -v $(pwd)/board.cfg:/app/instance/board.cfg gitlab-helper:<tag-version>
   ```
4 Waiting for process being taken util finished to exit the running container successfully.

5 Then you would be visit Gitlab service deployed by this helper.