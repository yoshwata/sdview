// process.stdin.resume();
// process.stdin.setEncoding('utf8');

// var inputText = "";
// process.stdin.on('data', function(chunk){
//   inputText += chunk;
// });

// process.stdin.on('end', function(){
//   inputText.split('\n').forEach(text => {
// 		console.log(text);
// 	});
// });

const axios = require('axios');
const config = require('./config');
axios.defaults.baseURL = config.baseURL;
const exec = async () => {
  try {
    const r = await axios.get(`/v4/auth/token?api_token=${config.usertoken}`)
    console.log(r.data);
    const jwt = r.data.token;
  } catch (error) {
    console.log(error);
  }

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

exec();


// axios.defaults.headers.common['Authorization'] = AUTH_TOKEN;
// axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
