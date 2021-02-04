import requests, logging, http, time, json, pprint
import service.config

log = logging.getLogger(__name__)
max_retries = 200

gitlab = service.config.get_config_from_file("gitlab")
gitlab_url = f'http://{gitlab["host_ip"]}:{gitlab["host_port"]}'
base_api_url = "{}/api/v4".format(gitlab_url)

def ping_gitlab():
  def request_with_status():
    try:
      resp = requests.get(url=gitlab_url)
      log.info("Requested Gitlab: %s with response status code: %s", gitlab_url, resp.status_code)
      return resp.status_code
    except Exception as e:
      log.error("Failed to request Gitlab with url: %s, error: %s", gitlab_url, e)
      return http.HTTPStatus.BAD_GATEWAY
  retries = 0
  while retries <= max_retries and request_with_status() >= http.HTTPStatus.BAD_REQUEST:
    retries += 1
    log.info("Retry to request Gitlab: %s for %d time(s)...", gitlab_url, retries)
    time.sleep(3)
  log.info("Gitlab service has been started successfully.")
  if retries > max_retries:
    raise("Failed to request Gitlab as exceeding max retries.")
  
def request_gitlab_api(method, token, api_url, **params):
  ping_gitlab()
  resp = requests.request(method=method, url="{}/{}".format(base_api_url, api_url), headers={"PRIVATE-TOKEN": token}, params=params)
  try:
    json_data = resp.json()
    log.info("Requested with URL: %s, status: %d, response: %s", api_url, resp.status_code, pprint.pprint(json_data))
    return json_data
  except json.JSONDecodeError as e:
    log.error("Failed to decode JSON: %s", e)
    return None
  
def get_application_settings(token):
  request_gitlab_api("GET", token, "application/settings")

def allow_local_request_webhooks(token):
  request_gitlab_api("PUT", token, "application/settings?{}={}".format("allow_local_requests_from_web_hooks_and_services", "true"))

def get_shared_runners(token):
  return request_gitlab_api("GET", token, "runners/all")

def delete_shared_runners(token, runner_id):
  log.info("Deleting shared runner with ID: %d", runner_id)
  request_gitlab_api("DELETE", token, "runners/{}".format(runner_id))