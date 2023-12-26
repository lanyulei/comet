import request from '/@/utils/request';

/**
 * 任务历史列表
 */
export function getTasks(params: object) {
	return request({
		url: `/apis/tekton.dev/v1/tasks`,
		method: 'get',
		params: params,
	});
}
