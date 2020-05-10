<template>
  <a-layout>
    <a-layout-header>
      <a-select
        size="large"
        mode="multiple"
        style="width: 100%"
        allowClear
        showArrow
        label-in-value
        :value="selectedTopics"
        :loading="loading"
        placeholder="Select topics"
        @change="handleTopicChange"
      >
        <a-spin v-if="loading" slot="notFoundContent" size="small" />
        <a-select-option v-for="item in topics" :key="item" :value="item">
          {{ item }}
        </a-select-option>
      </a-select>
    </a-layout-header>

      <a-layout-content>
        <div style="background-color: #ececec; padding: 20px;">
        <a-row :gutter="16">
          <a-col :span="6"  v-for="d in data" >
        <a-card  style="margin-bottom: 20px;" :headStyle="{'font-weight':'bold'}">
          <a slot="extra" v-if="d.type !== ''"><a-tag color="pink">{{d.data.type}}</a-tag></a>
          <a slot="title" :href="d.data.home_addr">{{d.data.title}}</a>
          <div slot="actions" :href="d.data.home_addr" v-if="topicsDiff(d.data.topics).length > 0">
            <a-button-group>
            <a-button type="dashed" v-for="t in topicsDiff(d.data.topics)" @click="addTopic(t)">
              {{t}}
            </a-button>
            </a-button-group>
          </div>
          <p>{{d.data.description}}</p>
        </a-card>

          </a-col>
        </a-row>
        </div>

      </a-layout-content>
    <a-layout-footer>footer</a-layout-footer>
  </a-layout>
</template>
<script>
    export default {
        data() {
            return {
                data: [],
                topics: [],
                selectedTopics: [],
                loading: false
            };
        },

        beforeMount() {
            this.loadData()
        },

        methods: {
            async loadData() {
                this.loading = true
                // todo: there must be a better way to searialize array into query string
                let path = '/entities/'
                let qs = []
                if (this.selectedTopics.length > 0) {
                    for (let i = 0; i < this.selectedTopics.length; i++) {
                        qs[i] = "topics=" + this.selectedTopics[i].key;
                    }
                }
                if (qs.length > 0) {
                    path += "?" + qs.join("&")
                }

                try {
                    await this.$axios.$get(path).then((resp) => {
                        this.data = resp.entities
                        this.topics = resp.topics
                    })
                } catch (e) {
                    this.$notify({
                        type: 'error',
                        title: 'Failed to load data',
                        text: e.response.data.error
                    })
                }
                this.loading = false
            },

            handleTopicChange(selectedTopics) {
                this.topics = []
                this.selectedTopics = selectedTopics;
                this.loadData()
            },

            topicsDiff(topics) {
                let sTopics = []
                for (let i = 0; i < this.selectedTopics.length; i++) {
                    sTopics.push(this.selectedTopics[i].key)
                }
                return topics.filter(x => !sTopics.includes(x));
            },

            addTopic(t) {
                this.selectedTopics.push({key: t, label: t})
                this.loadData()
            },
        },

        computed: {
        },
    };
</script>
