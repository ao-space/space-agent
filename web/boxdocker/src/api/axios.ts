/*
 * Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import axios from 'axios';

axios.defaults.headers.post['Content-Type'] = 'multipart/form-data';
axios.defaults.transformRequest = function (data) {
    let code = new FormData();
    for (let key in data) {
        code.append(key, data[key])
    }
    return code;
};

export function getAgentInfo() {
    return axios.get("/agent/info");
}

export function validateCode(tryoutCode,email) {
    return axios.post("/agent/v1/api/pair/tryout/code",{tryoutCode,email});
}

