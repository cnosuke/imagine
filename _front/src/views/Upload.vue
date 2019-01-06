<template>
  <div class="upload-container">
    <el-upload
      class="upload-demo"
      drag
      action=""
      :auto-upload="false"
      :on-change="onChange"
      :http-request="httpRequest"
      ref="upload"
      multiple
    >
      <i class="el-icon-upload"></i>
      <div class="el-upload__text">
        ここにファイルをドラッグ&ドロップしてください<br />または、<em
          >ここをクリックしてファイルを選択</em
        >
      </div>
      <div class="el-upload__tip" slot="tip">
        1GiB以上はアップロードできません
      </div>
    </el-upload>

    <el-button
      type="success"
      @click="submitUpload"
      :disabled="submitButtonDisable()"
      :loading="isLoading()"
      class="upload-button"
      >upload to server</el-button
    >

    <div class="uploaded-list" v-if="doneList.length > 0">
      <h4>Uploaded List</h4>
      <el-table :data="doneList" style="width: 100%">
        <el-table-column label="FileName" prop="filename"> </el-table-column>
        <el-table-column label="URL" prop="url"> </el-table-column>
        <el-table-column align="right">
          <template slot-scope="scope">
            <el-button
              size="mini"
              type="plain"
              @click="copyUrl(scope.$index, doneList)"
              >URL Copy</el-button
            >
            <el-button
              size="mini"
              type="primary"
              @click="openUrl(scope.$index, doneList)"
              >Open</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  ElUploadInternalFileDetail,
  ElUploadProgressEvent,
  HttpRequestOptions,
} from "element-ui/types/upload";
import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
import mime from "mime";
import { Clipboard } from "ts-clipboard";

@Component({
  components: {},
})
export default class extends Vue {
  public fileList: ElUploadInternalFileDetail[] = [];
  public doneList: Array<{ filename: string; s3Key: string; url: string }> = [];

  constructor() {
    super();
  }

  private onChange(
    f: ElUploadInternalFileDetail,
    fList: ElUploadInternalFileDetail[],
  ) {
    this.fileList = fList;
  }

  private submitUpload() {
    (this.$refs.upload as any).submit();
  }

  private httpRequest(req: HttpRequestOptions) {
    const f = req.file;
    const t = this.fetchContentType(f.name);

    axios
      .post("//localhost:8888/api/v1/create_presigned_post_url", {
        filename: f.name,
        content_type: t,
      })
      .then((res: AxiosResponse) => {
        const presignedPostUrl = res.data.url;
        const reqConfig: AxiosRequestConfig = {
          onUploadProgress: (e: ElUploadProgressEvent): void => {
            e.percent = Math.round((e.loaded * 100) / e.total);
            req.onProgress(e);
          },
          headers: { "Content-Type": res.data.content_type },
        };

        axios
          .put(presignedPostUrl, f, reqConfig)
          .then((res0: AxiosResponse) => {
            this.doneList.push({
              filename: res.data.filename as string,
              s3Key: res.data.key as string,
              url: this.buildDownloadUrl(res.data.key as string),
            });

            req.onSuccess(res0);
          })
          .catch((err0: ErrorEvent) => {
            req.onError(err0);
          });
      })
      .catch((err: ErrorEvent) => {
        req.onError(err);
      });
  }

  private submitButtonDisable(): boolean {
    return this.fileList.length === 0 || this.isLoading();
  }

  private isLoading(): boolean {
    let fl: boolean = false;
    this.fileList.forEach((f: ElUploadInternalFileDetail) => {
      if (f.status === "uploading") {
        fl = true;
        return;
      }
    });

    return fl;
  }

  private copyUrl(
    idx: number,
    list: Array<{ filename: string; s3Key: string; url: string }>,
  ) {
    Clipboard.copy(list[idx].url);
    this.$message(`${list[idx].filename} のURLをコピーしました`);
  }

  private openUrl(
    idx: number,
    list: Array<{ filename: string; s3Key: string; url: string }>,
  ) {
    window.open(list[idx].url, "_blank");
  }

  private fetchContentType(filename: string): string {
    const n = filename.toLowerCase();
    const ext = n.split(".").pop();
    if (ext) {
      return mime.getType(ext) || "application/octet-stream";
    } else {
      return "application/octet-stream";
    }
  }

  private buildDownloadUrl(key: string): string {
    return `https://example.com/${key}`;
  }
}
</script>

<style lang="less" scoped>
.upload-container {
  width: 800px;
  margin: 0 auto;
}

.upload-button {
  margin-top: 20px;
}

.uploaded-list {
  margin-top: 60px;
}
</style>
