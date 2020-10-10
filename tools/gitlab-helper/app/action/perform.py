from service.ssh import SSHUtil
import service.config
import service.http
import logging
import re, time
from os import path
import sys, getopt

log = logging.getLogger(__name__)

def detect_gitlab_status():
  gitlab = service.config.get_config_from_file("gitlab")
  command_detection = f'''if [ -z $(docker ps -f name={gitlab['container_name']} --format {{{{.Names}}}}) ] && \
[ ! -z $(docker ps -a -f name={gitlab['container_name']} --format {{{{.Names}}}}) ]; then \
echo "Removing stopped Gitlab container...";\
docker rm {gitlab['container_name']};\
fi'''
  return SSHUtil.exec_command(command_detection)

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
  gitlab = service.config.get_config_from_file("gitlab")
  return f'''docker exec -i {gitlab["container_name"]} bash gitlab-rails runner {command_line}'''

def reset_root_password():
  log.info("Resetting root password ...")
  gitlab = service.config.get_config_from_file("gitlab")
  command_reset_root_password=f'''\'user = User.find_by(username: "root"); user.password = "{gitlab["root_password"]}"; user.password_confirmation = "{gitlab["root_password"]}"; user.save!\''''
  return SSHUtil.exec_command(gitlab_docker_exec(command_reset_root_password))

def setting_access_token(token):
  log.info("Setting root access token ...")
  access_token_name = "auto_gen_token"
  command_set_access_token=f'''\'user = User.where(admin: true).first; user.personal_access_tokens.where("name":"{access_token_name}").each{{|t| t.delete}}; token = user.personal_access_tokens.create(scopes: [:read_user, :read_repository, :write_repository, :api, :sudo], name: "{access_token_name}"); token.set_token("{token}"); token.save!\''''
  return SSHUtil.exec_command(gitlab_docker_exec(command_set_access_token))

def update_access_token(token):
  try:
    config_file_path = path.join(path.dirname(path.dirname(path.abspath(__file__))), "instance", "board.cfg")
    log.debug("Config file path: %s with access_token: %s", config_file_path, token)
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

def obtain_shared_runner_token():
  log.info('Obtaining token for shared Gitlab runner...')
  cmd_obtain_shared_runner_token = f'''\'puts Gitlab::CurrentSettings.current_application_settings.runners_registration_token\''''
  token = SSHUtil.exec_command(gitlab_docker_exec(cmd_obtain_shared_runner_token))
  if token == "":
    log.error("Failed to obtain Gitlab shared runner.")
    return None
  return token.strip()

def reset_gitlab_runner_config():
  log.info("Reset Gitlab runner config.toml")
  default_config = '''concurrent = 1
check_interval = 0
[session_server]
session_timeout = 1800'''
  r = service.config.get_config_from_file("gitlab-runner")
  cmd_reset_config = ""
  for index, content in enumerate(default_config.splitlines()):
    append_operator = ">>"
    if index == 0:
      append_operator = ">"
    cmd_reset_config += f'''echo {content} {append_operator} {r["config_path"]};'''
  return SSHUtil.exec_command(cmd_reset_config)

def clean_up_stale_runner(token):
  log.info("Clean up stale Gitlab shared runners...")
  runner_list = service.http.get_shared_runners(token)
  for runner in runner_list:
    service.http.delete_shared_runners(token, int(runner["id"]))

def command_to_register_runner(gitlab_url, gitlab_runner_token, executor, r):
  cmd = f'''gitlab-runner register --name "{r["runner_name"] + "-" + executor}" \
--url="{gitlab_url}" \
--registration-token="{gitlab_runner_token}" \
--executor="{executor}" \
--non-interactive \
--tag-list "{executor + "-ci"}"'''
  if executor == "docker":
    cmd += f''' --docker-image "{r["runner_image"]}"'''
  return cmd

def register_gitlab_shared_runner(gitlab_runner_token, executor):
  log.info("Registering Gitlab runner with token: %s", gitlab_runner_token)
  gitlab = service.config.get_config_from_file("gitlab")
  gitlab_url = f'''http://{gitlab["host_ip"]}:{gitlab["host_port"]}'''
  r = service.config.get_config_from_file("gitlab-runner")
  cmd_runner_register = command_to_register_runner(gitlab_url, gitlab_runner_token, executor, r)
  return SSHUtil.exec_command(cmd_runner_register)

if __name__ == '__main__':
  logging.basicConfig(level=logging.INFO)
  try:      
    opts, args = getopt.getopt(sys.argv[1:], "hr:",["reset-token-only=",])
    reset_token_only = False
    for opt, arg in opts:
      if opt in ("-r", "--reset-token-only"):
        if arg and arg.lower() == "true":
          reset_token_only = True
    log.info("Obtaining Gitlab admin access token...")
    admin_access_token = service.config.generate_token()
    if reset_token_only:
      log.info("Resetting token only...")
      setting_access_token(admin_access_token)
      update_access_token(admin_access_token)
    else:
      log.info("Start normally...")
      detect_gitlab_status()
      gitlab_docker_run()
      setting_access_token(admin_access_token)
      update_access_token(admin_access_token)
      update_allow_local_webhook_request(admin_access_token)
      runner_token = obtain_shared_runner_token()
      if runner_token:
        reset_gitlab_runner_config()
        clean_up_stale_runner(admin_access_token)
        register_gitlab_shared_runner(runner_token, "docker")
        register_gitlab_shared_runner(runner_token, "shell")
  except getopt.GetoptError:
    log.info("action/perform.py -ro | --reset-token-only=[true]")