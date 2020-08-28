from service.ssh import SSHUtil
import service.config
import service.http
import logging
import re
from os import path

log = logging.getLogger(__name__)

def gitlab_docker_run():
  gitlab = service.config.get_config_from_file("gitlab")
  command_gitlab_run = f'''docker run -d \
-p {gitlab["host_port"]}:{gitlab["container_port"]} \
-p {gitlab["host_ssh_port"]}:{gitlab["container_ssh_port"]} \
-v {gitlab["custom_config"]}:/etc/gitlab/gitlab.rb \
-v {gitlab["base_data_path"]}/config:/etc/gitlab:Z \
-v {gitlab["base_data_path"]}/logs:/var/log/gitlab:Z \
-v {gitlab["base_data_path"]}/data:/var/opt/gitlab:Z \
--name {gitlab["container_name"]} {gitlab["image_name"]}'''
  return SSHUtil.exec_command(command_gitlab_run)

def gitlab_docker_exec(command_line):
  service.http.ping_gitlab()
  return "docker exec -i gitlab bash gitlab-rails runner {}".format(command_line)

def reset_root_password():
  log.info("Resetting root password ...")
  gitlab = service.config.get_config_from_file("gitlab")
  command_reset_root_password=f'''\'user = User.find_by(username: "root"); user.password = "{gitlab["root_password"]}"; user.password_confirmation = "{gitlab["root_password"]}"; user.save!\''''
  return SSHUtil.exec_command(gitlab_docker_exec(command_reset_root_password))

def setting_access_token(token):
  log.info("Setting root access token ...")
  command_set_access_token=f'''\'user = User.where(id: 1).first; token = user.personal_access_tokens.create(scopes: [:read_user, :read_repository, :write_repository, :api, :sudo], name: "Automation_token_updated"); token.set_token("{token}"); token.save!\''''
  return SSHUtil.exec_command(gitlab_docker_exec(command_set_access_token))

def update_access_token(token):
  try:
    config_file_path = path.join(path.dirname(path.dirname(path.abspath(__file__))), "instance", "board.cfg")
    log.debug("Config file path: %s", config_file_path)
    with open(config_file_path, "r") as f:
      content = f.read()
      content_updates = re.sub(r"^(gitlab_admin_token\s*=\s*)(.*)$", r"\g<1>{}".format(token), content, flags=re.M)
    with open(config_file_path, "w") as f:
      f.write(content_updates)
    log.info("Successful updated root access token to board.cfg.")
  except Exception as e:
    log.error("Failed to update config file: %s", e)

def update_allow_local_webhook_request(token):
  service.http.allow_local_request_webhooks(token)

def get_application_settings(token):
  service.http.get_application_settings(token)

if __name__ == '__main__':
  logging.basicConfig(level=logging.INFO)
  log.info(gitlab_docker_run())
  # log.info(reset_root_password())
  
  admin_access_token = service.config.generate_token()
  log.info(setting_access_token(admin_access_token))
  update_access_token(admin_access_token)
  update_allow_local_webhook_request(admin_access_token)
  get_application_settings(admin_access_token)