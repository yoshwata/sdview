

const axios = require('axios');
const config = require('./config');
axios.defaults.baseURL = config.baseURL;
const exec = async ({ buildId }) => {
  try {
    const r = await axios.get(`/v4/auth/token?api_token=${config.usertoken}`)
    // console.log(r.data);
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

    console.log(b.data.meta.build.pipelineId);
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

  // apiでrunnningビルドを全部取るみたいなことはできないっすね
  // いいのじゃ、kubectlでとるから
  // .then(function (response) {
  //   // handle success
  //   console.log(response);
  // })
  // .catch(function (error) {
  //   // handle error
  //   console.log(error);
  // })
  // .then(function () {
  //   // always executed
  // });
}

// axios.defaults.headers.common['Authorization'] = AUTH_TOKEN;
// axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';

// process.stdin.resume();
// process.stdin.setEncoding('utf8');

var inputText = "";
process.stdin.on('data', function (chunk) {
  inputText += chunk;
});

process.stdin.on('end', function () {
  inputText.split('\n').forEach(text => {
    console.log("aaaaaaaaaaaaaa" + text)
    if (text) {
      exec({ buildId: text });
    }
  });
});
