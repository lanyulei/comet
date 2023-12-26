import { createApp } from 'vue';
import pinia from '/@/stores/index';
import App from '/@/App.vue';
import router from '/@/router';
import { directive } from '/@/directive';
import { i18n } from '/@/i18n';
import other from '/@/utils/other';

import ElementPlus from 'element-plus';
import '/@/theme/index.scss';
import VueGridLayout from 'vue-grid-layout';
// 图标统一注册
import * as icons from '@element-plus/icons-vue';

const app = createApp(App);

Object.keys(icons).forEach((key) => {
	app.component(key, (icons as any)[key]);
});

directive(app);
other.elSvg(app);

app.use(pinia).use(router).use(ElementPlus).use(i18n).use(VueGridLayout).mount('#app');
