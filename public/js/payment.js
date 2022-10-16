(self["webpackChunkvue_boilerplate"] = self["webpackChunkvue_boilerplate"] || []).push([["/js/payment"],{

/***/ "./resources/src/js/payment.js":
/*!*************************************!*\
  !*** ./resources/src/js/payment.js ***!
  \*************************************/
/***/ (function() {

paypal.Buttons({
  // Set your environment
  env: paypal_env,
  // Set style of buttons
  style: {
    layout: 'horizontal',
    // horizontal | vertical
    size: 'large',
    // medium | large | responsive
    shape: 'rect',
    // pill | rect
    color: 'gold',
    // gold | blue | silver | black,
    fundingicons: false,
    // true | false,
    tagline: false // true | false,
  },

  // onInit is called when the button first renders
  onInit: function onInit(data, actions) {
    // Disable the buttons
    actions.disable();

    // Listen for changes to the checkbox
    document.querySelector('#amount').addEventListener('keyup', function (event) {
      var amount = event.target.value;
      amount = +amount;
      if (isNaN(amount)) {
        document.getElementById("amount-err").innerHTML = "Invalid amount format";
        setTimeout(function () {
          document.getElementById("amount-err").innerHTML = "";
        }, 2000);
        actions.disable();
      } else if (amount < 5) {
        document.getElementById("amount-err").innerHTML = "Amount must be greater than 5";
        setTimeout(function () {
          document.getElementById("amount-err").innerHTML = "";
        }, 2000);
        actions.disable();
      }
      if (event.target.value === "") {
        document.getElementById("amount-err").innerHTML = "Invalid amount format";
        setTimeout(function () {
          document.getElementById("amount-err").innerHTML = "";
        }, 2000);
        actions.disable();
      } else {
        actions.enable();
      }
    });
  },
  // onClick is called when the button is clicked
  onClick: function onClick() {
    var regex = /^\d+(?:\.\d{0,2})$/;
    var amount = document.getElementById("amount").value;
    amount = +amount;
    if (isNaN(amount)) {
      alert(2);
      document.getElementById("amount-err").innerHTML = "Invalid amount format";
      setTimeout(function () {
        document.getElementById("amount-err").innerHTML = "";
      }, 2000);
    }
  },
  // Wait for the PayPal button to be clicked
  createOrder: function createOrder() {
    var formData = new FormData();
    var amount = document.getElementById("amount").value;
    formData.append('amount', amount);
    return fetch('/app/paypal/do/order', {
      method: 'POST',
      body: formData
    }).then(function (response) {
      return response.json();
    }).then(function (resJson) {
      return resJson.data.id;
    });
  },
  // Wait for the payment to be authorized by the customer
  onApprove: function onApprove(data, actions) {
    return fetch('/app/paypal/order/success/' + data.orderID, {
      method: 'POST'
    }).then(function (res) {
      return res.json();
    }).then(function (res) {

      // window.location.href = '/account/paypal/order/success';
    });
  },
  // Wait for the payment to be authorized by the customer
  onCancel: function onCancel(data, actions) {
    return fetch('/app/paypal/order/cancel/' + data.orderID, {
      method: 'POST'
    }).then(function (res) {
      return res.json();
    }).then(function (res) {

      // window.location.href = '/account/paypal/order/success';
    });
  }
}).render('#paypalCheckoutContainer');

/***/ })

},
/******/ function(__webpack_require__) { // webpackRuntimeModules
/******/ var __webpack_exec__ = function(moduleId) { return __webpack_require__(__webpack_require__.s = moduleId); }
/******/ var __webpack_exports__ = (__webpack_exec__("./resources/src/js/payment.js"));
/******/ }
]);
//# sourceMappingURL=payment.js.map