from paramiko import SSHClient, AutoAddPolicy
import service.config

import logging, sys

log = logging.getLogger(__name__)

class SSHUtil:
  client = None
  @classmethod
  def __get_conn(cls):
    cls.client = SSHClient()
    cls.client.set_missing_host_key_policy(AutoAddPolicy())
    gitlab = service.config.get_config_from_file("gitlab")
    ssh = service.config.get_config_from_file("ssh")
    for item in ["ssh_username", "ssh_password"]:
      if ssh[item] is None:
        raise("Missing config of {}, aborting...")
    cls.client.connect(hostname=gitlab["host_ip"], port=ssh["ssh_port"], username=ssh["ssh_username"], password=ssh["ssh_password"])
    
  @classmethod
  def exec_command(cls, command, *args):
    try:
      cls.__get_conn()
      command = '{} {}'.format(command, ' '.join(args))
      log.info("Executing command: {}".format(command))
      _, stdout, stderr = cls.client.exec_command(command)
      # if stderr:
      #   return stderr.read().decode('utf-8')
      return stdout.read().decode('utf-8')
    except Exception as e:
      log.error("Failed to execute command via SSH: %s", e)
      sys.exit(255)
    finally:
      cls.client.close()