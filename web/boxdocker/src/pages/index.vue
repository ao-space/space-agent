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
      <div class="container flex-column flex-xy-center" v-loading="loading">
        <div class="black-20 fw-b ">傲空间私有部署激活</div>
        <div style="position: relative;margin-top:68px;">
          <div class="input-box black-16">
              <div class="flex-y-center">
                <label>邮箱</label>
                <input ref="firstinput" class="black-16" v-model="email" type="text" placeholder="请输入您申请时的邮箱"/>
              </div>
              <div class="flex-y-center" style="border-top:1px solid #BCBFCD;">
                <label>激活码</label>
                <input class="black-16" v-model="code" type="text" placeholder="请输入激活码"/>
              </div>
          </div>
          <div v-if="error" class="error red-14">{{ error }}</div>
          <div class="mt-50 black56-14 flex-y-center">
            <img src="@/assets/svg/ts.svg" />
            <span class="ml-5 pointer" @click="show = !show"
              >如何获取激活码</span
            >
          </div>
        </div>
        <div :class="{'button-blue':commit,'button-gray':!commit}" @click="checkCode">提交</div>
      </div>
    </div>
    <div class="mt-40 black56-14">
      copyright © 2022-2023 中国科学院软件研究所
    </div>
    <el-dialog v-model="show">
      <div class="mt-30 mb-20 tc black-16 fw-b">如何获取激活码</div>
      <ul class="list black-14">
        <li>傲空间官网 ao.space 或 关注傲空间公众号点击 “加入公测” ，填写您的邮箱，审核通过后会将激活码发送到您的邮箱，数量有限，先到先得。</li>
        <li>激活码有效期 24 小时，失效后请重新获取激活码。</li>
        <li>每个激活码仅可使用一次，使用后即失效。</li>
      </ul>
      <div class="flex-x-center mt-30">
        <div class="img mr-40">
          <img src="@/assets/png/gw@2x.png" />
          <div class="apply">官网申请</div>
        </div>
        <div class="img">
          <img src="@/assets/png/gzh@2x.png" />
          <div class="apply">公众号申请</div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>

import {validateCode} from "../api/axios";
import { getAgentInfo } from "@/api/axios";

export default {
  name: "Index",
  components: {
  },
  computed:{
    commit(){
      return this.email.length > 0 && this.code.length>0
    }
  },
  data() {
    return {
      error: '',
      show: false,
      code: '',
      email: '',
      loading: true,
    };
  },
  mounted() {
    getAgentInfo().then((result) => {
      if (result.data.code === "AG-200") {
        let tryoutCodeVerified = result.data.results.tryoutCodeVerified
        if(tryoutCodeVerified){
          this.$router.push("/code");
        }
      }
    }).finally(()=>{
      this.loading = false
      this.$nextTick(() => {
        this.$refs.firstinput.focus()
      })
    })
  },
  methods: {
    checkCode() {
      if(!this.commit){
        return
      }
      if (this.isNotEmail()) {
        this.error = '邮箱格式不正确'
        return
      }

      validateCode(this.code,this.email).then((result)=>{
        console.log(result)
        if(result.data.code === 'AG-200') {
          this.$router.push("/code");
        }else if(result.data.code === 'AG-465'){
          this.error = '激活码错误，请确认后重新输入(AG-465)'
        }else if(result.data.code === 'AG-466'){
          this.error = '激活码已失效，请重新获取(AG-466)'
        }else if(result.data.code === 'AG-467'){
          this.error = '该邮箱已注册，请更换邮箱重新申请(AG-467)'
        }else if(result.data.code === 'AG-468'){
          this.error = '激活码无效，请重新获取(AG-468)'
        }else if(result.data.code === 'AG-469'){
          this.error = '容器下载中，请稍后重试(AG-469)'
        }else{
          this.error = '系统服务错误，请稍后重试'
        }
      })
    },
    isNotEmail() {
      let emailPat = /^(.+)@(.+)$/
      let result = this.email.match(emailPat)
      return result === null
    },
  },
};
</script>

<style scoped lang="scss">
::v-deep(.el-loading-mask){
  border-radius: 10px;
}
.input-box{
  width: 500px;
  height: 100px;
  background: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #BCBFCD;
  div{
    height: 50px;
    padding-left: 20px;
    label{
      width: 68px;
    }
    input{
      width: 390px;
      border: none;
      outline: none;
    }
  }
}
.img{
  position: relative;
  width: 262px;
  height: 212px;
  img{
    width: 100%;
    height: 100%;
  }
  .apply{
    position: absolute;
    left: 10px;
    bottom: 10px;
    font-size: 16px;
    font-weight: bolder;
    color: #FFFFFF;
  }
}
.list{
  margin-left: 122px;
  list-style: decimal;
  line-height: 22px;
}
.error {
  position: absolute;
  top: 117px;
}
.button-blue {
  width: 250px;
  height: 50px;
  margin-top: 60px;
  font-size: 20px;
}
.button-gray {
  color: #FFFFFF;
  background: #C8C8C8;
  box-shadow: 0 0 3px 0 #C0C9D8;
  border-radius: 10px;
  @extend .button-blue
}
.container {
  width: 900px;
  height: 528px;
  border-radius: 10px;
  backdrop-filter: blur(2px);
  background: rgba(255, 255, 255, 0.6);
}
::v-deep(.el-dialog__header) {
  border: none;
}
::v-deep(.el-dialog) {
  width: 800px;
  height: 510px;
  background: #ffffff;
  border-radius: 6px;
  --el-dialog-margin-top: 24vh;
}
.content {
  max-width: 1920px;
  margin-top: 170px;
}
</style>
