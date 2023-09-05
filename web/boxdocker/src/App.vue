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
  <el-container>
    <el-header>
      <div class="header flex-y-center">
        <div class=" m-center flex-y-center" style="width: 1920px">
          <img class="logo" src="@/assets/png/logo@2x.png" />
        </div>
      </div>
    </el-header>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
</template>
<style lang="scss" scoped>
.header{
  height: 60px;
  background: rgba(255,255,255,0.6);
}
.logo {
  width: 126px;
  height: 46px;
  margin-left: 26px;
}
.el-divider--vertical{
  margin-left: 0;
  margin-right: 20px;
}
</style>
