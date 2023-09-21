<!--
  ~ Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

<template>
  <div class="content m-center flex-xy-center flex-column">
    <div>
      <div class="container flex-column flex-y-center">
        <div class="black-20 fw-b mt-44">傲空间设备码</div>
        <img class="img" :src="imgUrl" />
        <div class="black-18 fw-b">
          请使用 <span class="blue-18">傲空间开源版 App</span> 扫码绑定后使用
        </div>
        <div class="mt-20 gray-14">
          未安装傲空间 App 请先 <a target="_blank" href="https://ao.space/download/opensource" class="download blue-14">下载傲空间</a>
        </div>
      </div>
    </div>
    <div class="mt-40 black56-14">
      copyright © 2022-2023 中国科学院软件研究所
    </div>
  </div>
</template>

<script>
import { getAgentInfo } from "@/api/axios";
import QRCode from "qrcode";

export default {
  name: "Code",
  data() {
    return {
      imgUrl: "",
    };
  },
  mounted() {
    getAgentInfo().then((result) => {
      if (result.data.code === "AG-200") {
        let url = result.data.results.boxQrCode;
        QRCode.toDataURL(
          url,
          { errorCorrectionLevel: "L", margin: 2, width: 220 },
          (err, url) => {
            this.imgUrl = url;
          }
        );
      }
    });
  },
};
</script>

<style scoped lang="scss">
.img {
  width: 220px;
  height: 220px;
  background: #ffffff;
  margin: 50px;
}
.download {
  text-decoration: underline;
}
.container {
  width: 900px;
  height: 500px;
  background: rgba(255,255,255,0.6);
  border-radius: 10px;
  backdrop-filter: blur(2px);
}
.content {
  max-width: 1920px;
  margin-top: 170px;
}
</style>
