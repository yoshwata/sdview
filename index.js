#! /usr/bin/env node

const axios = require('axios');
const config = require('./config');
axios.defaults.baseURL = config.baseURL;
const exec = async ({ buildId }) => {
  try {
    const r = await axios.get(`/v4/auth/token?api_token=${config.usertoken}`)
    const jwt = r.data.token;

    const b = await axios.get(
      `/v4/builds/${buildId}`,
      {
        headers: {
          "accept": "application/json",
          "Authorization": jwt
        }
      }
    );

    const pipelineId = b.data.meta.build.pipelineId;

    const p = await axios.get(
      `/v4/pipelines/${pipelineId}`,
      {
        headers: {
          "accept": "application/json",
          "Authorization": jwt
        }
      }
    );

    const csv = buildId + ',' + pipelineId + ',' + p.data.name;

    console.log(csv);

  } catch (error) {
    console.log(error);
  }
}

var inputText = "";
process.stdin.on('data', function (chunk) {
  inputText += chunk;
});

process.stdin.on('end', function () {
  inputText.split('\n').forEach(text => {
    if (text) {
      exec({ buildId: text });
    }
  });
});
