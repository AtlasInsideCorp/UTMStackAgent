import argparse
import json
import os
import requests
import pytz
from datetime import datetime
from google.cloud import storage
from dotenv import load_dotenv

def main(update_version):
	load_dotenv()

	gcp_key = json.loads(os.environ.get("GCP_KEY"))

	# os.environ["GOOGLE_APPLICATION_CREDENTIALS"] = r'utmstack-377415-dbe8b897390f.json'
	storage_client = storage.Client.from_service_account_info(gcp_key)

    # Download version.json
	us_tz = pytz.timezone('America/New_York')
	us_time = datetime.now(us_tz)
	us_time_ms = int(us_time.timestamp()) * 1000

	version_json_url = "https://storage.googleapis.com/utmstack-updates/agent_updates/versions.json?time="+str(us_time_ms)
	version_data = json.loads(requests.get(version_json_url).content.decode('utf-8'))

	bucket_name = "utmstack-updates"
	bucket = storage_client.bucket(bucket_name)

	if update_version.startswith("dev-") or update_version.startswith("alpha-"):
		environment, version_number = update_version.split("-")
	elif update_version.startswith("v"):
		environment, version_number = "release", update_version[1:]
	else:
		raise ValueError("unknown environment")
	
	if environment == "dev":
		windows_key = "testing_windows"
		linux_key = "testing_linux"
		version_key = "testing_version"
	elif environment == "alpha":
		windows_key = "alpha_windows"
		linux_key = "alpha_linux"
		version_key = "alpha_version"
	elif environment == "release":
		windows_key = "windows_agent"
		linux_key = "linux_agent"
		version_key = "agent_version"
    

    # Upload compiled files
	windows_blob = bucket.blob(version_data[windows_key])
	windows_blob.upload_from_filename(os.path.join(os.environ["GITHUB_WORKSPACE"], "utmstack-windows.exe"))

	linux_blob = bucket.blob(version_data[linux_key])
	linux_blob.upload_from_filename(os.path.join(os.environ["GITHUB_WORKSPACE"],"utmstack-linux"))

    # Update version.json
	version_data[version_key] = version_number
	version_json_path = "agent_updates/versions.json"
	version_blob = bucket.blob(version_json_path)
	version_blob.upload_from_string(json.dumps(version_data))

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Update UTMStack agents in Google Cloud Storage")
    parser.add_argument("update_version", type=str, help="Update version string (e.g. dev-20060323, alpha-20060323, v1.0.1)")
    args = parser.parse_args()
    main(args.update_version)