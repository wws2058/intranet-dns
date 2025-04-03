import { message } from "ant-design-vue";
const success = () => {
  message.success("This is a success message");
};
const error = () => {
  message.error("This is an error message");
};
const warning = () => {
  message.warning("This is a warning message");
};

function successWithMsg(tips) {
  message.success(tips);
}

function myPopMessage() {}

export { success, error, warning, successWithMsg, myPopMessage };
