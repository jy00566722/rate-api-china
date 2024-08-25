export default defineBackground(() => {
  console.log('Hello background!', { id: chrome.runtime.id });
  // See https://wxt.dev/api/background.html
  const getData = async () => {
    const res = await fetch('https://api-test.jooop.top/v1/rate_usa');
    const data = await res.json();
    if(data && data.code == 2000){
      console.log(data)
      chrome.storage.local.set(data.data,()=>{
        console.log("set data success");
        //保存当前时间的本地UTC字符串到storage
        chrome.storage.local.set({lastUpdateTime: new Date().toUTCString()},()=>{
          console.log("set lastUpdateTime success");
        });
      }); 
    }
  }
  getData();
  //设置alarm每5分钟执行一次
  chrome.alarms.create("getData", { periodInMinutes: 2 });
  //监听alarm
  chrome.alarms.onAlarm.addListener((alarm) => {
    if(alarm.name == "getData"){
      getData();
    }
  });

});
