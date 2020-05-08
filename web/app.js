// ./app.js

var host = "http://localhost:5000"

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
      var c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }


new Vue({
    el: '#app',
    data: {
        resultAlbums: [],
    },
    methods: {        
        queryAlbums(){
            console.log("enter queryAlbums()")
            console.log("Start to list albums")
            axios.get(host + "/albums?sessionID="+getCookie("goquestsession"))
                .then(response => { 
                    console.log(response)
                    this.resultAlbums = response.data;
                })
            console.log("data")    
            console.log(this.resultAlbums)
        }
    }
});
