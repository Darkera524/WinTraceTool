########################################################################
# Copyright 2017 FireEye Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
########################################################################

import time
import etw
import requests
import sys
import json


def some_func(name, guid):
    # define capture provider info "{11111111-1111-1111-1111-111111111111}"
    providers = [etw.ProviderInfo(name, etw.GUID(guid))]
    # create instance of ETW class
    job = etw.ETW(providers=providers, event_callback=lambda x: print(str(x).replace("'","\"")))

    # start capture
    job.start()

    # wait some time
    #time.sleep(5)

    while True:
        url = "http://127.0.0.1:8093/query"
        d = [
            {
                "Provider": guid
            }
        ]

        r = requests.post(url, json.dumps(d))
        response = r.text

        if response == "no":
            # stop capture
            job.stop()
            break


if __name__ == '__main__':
    name = sys.argv[1]
    guid = sys.argv[2]
    some_func(name, guid)
