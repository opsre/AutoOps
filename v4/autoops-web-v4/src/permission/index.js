/**
 * vue3 自定义权限指令配置
 *
 * @author xiaoRui
 */
import authority from "./Authority"

export default {
    install(app) {
        app.directive('authority', authority)
    }
}
