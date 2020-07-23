from os import environ, path, SEEK_SET
import configparser as cp, re, logging, random, string
from io import StringIO

log = logging.getLogger(__name__)

def __resolve_file_path(sub_path, filename):
  file_path = path.join(path.dirname(path.dirname(__file__)), sub_path, filename)
  log.debug("Reading config from: %s", file_path)
  return file_path

def get_config_from_file(section="default"):
  try:
    c = cp.ConfigParser(allow_no_value=True)
    config_file_path = __resolve_file_path("env", "config.ini")
    c.read(config_file_path)
    c = __fusion_config(c)
    return c[section]
  except Exception as e:
    log.error("Failed read config file: %s", e)

def __get_config_from_cfg():
  try:
    w = StringIO()
    w.write("[board_cfg]\n")
    w.write(open(__resolve_file_path("instance", "board.cfg")).read())
    w.seek(0, SEEK_SET)
    c = cp.ConfigParser()
    c.read_file(w)
    return c
  except Exception as e:
    log.error("Failed to read board.cfg file: %s", e)
  return None

def __fusion_config(current): 
  updates = __get_config_from_cfg()
  if not updates:
    raise("Failed to load updates from board.cfg.")
  default_section = "board_cfg"
  for section in current.sections():
    for key in current.options(section):
      if updates.has_option(default_section, "gitlab_" + key):
        new_val = updates.get(default_section, "gitlab_" + key)
        current.set(section, key, new_val)
        log.info("Overwriting config: %s with value: %s in section: %s", key, new_val, section)
  return current

def generate_token(length=20):
  return ''.join(random.choice(string.ascii_letters) for i in range(length))
