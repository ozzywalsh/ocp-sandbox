import logging
import os
import random
import time

import requests

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
log = logging.getLogger("traffic-generator")

TARGET_URL = os.environ.get("TARGET_URL", "http://taxform-renderer.sandbox.svc:8080/render")
INTERVAL = int(os.environ.get("INTERVAL_SECONDS", "10"))
ERROR_RATE = float(os.environ.get("ERROR_RATE", "0.2"))

VALID_PAYLOAD = {
    "template": "foo",
    "fields": {
        "nome completo": "Test User",
        "data de nascimento": "15/06/1985",
        "nacionalidade": "Brasileira",
        "nome do pai": "Carlos Silva",
        "nome da mãe": "Maria Silva",
        "endereço completo rua número município  local país": "Rua Example 123, São Paulo, Brasil",
        "Local e Data": "São Paulo, 01/01/2026",
    },
    "radioButtonGroups": {
        "Condição": "residente",
    },
}


def send_request():
    if random.random() < ERROR_RATE:
        log.info("Sending invalid request (empty body)")
        resp = requests.post(TARGET_URL, timeout=30)
    else:
        log.info("Sending valid render request")
        resp = requests.post(TARGET_URL, json=VALID_PAYLOAD, timeout=30)

    log.info("Response: status=%d size=%d", resp.status_code, len(resp.content))


if __name__ == "__main__":
    log.info("Starting traffic generator targeting %s every %ds (error_rate=%.0f%%)",
             TARGET_URL, INTERVAL, ERROR_RATE * 100)
    while True:
        try:
            send_request()
        except Exception:
            log.exception("Request failed")
        time.sleep(INTERVAL)
