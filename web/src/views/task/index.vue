<template>
	<div class="main-container">
		<el-alert title="Task 是一个可重复使用的工作单元，定义了容器化的工作步骤，用于构建和部署应用程序。" class="mb10" type="info" show-icon />
		<el-card class="box-card">
			<div>
				<el-button type="primary" icon="ele-Plus" class="mr8">新建</el-button>
				<el-input v-model="listQuery.name" placeholder="请输入搜索项" @keyup.enter.native="getList" style="max-width: 300px">
					<template #append>
						<el-button icon="ele-Search" @click="getList"></el-button>
					</template>
				</el-input>
			</div>
			<div class="mt10">
				<el-table :data="list" border stripe style="width: 100%">
					<el-table-column prop="metadata.name" label="名称" />
					<el-table-column prop="metadata.namespace" label="命名空间" />
					<el-table-column label="Step 数量">
						<template #default="{ row }">
							<span>{{ row.spec?.steps?.length ? row.spec?.steps?.length : 0 }}</span>
						</template>
					</el-table-column>
					<el-table-column prop="metadata.creationTimestamp" label="创建时间">
						<template #default="{ row }">
							<span>{{ parseTime(row.metadata.creationTimestamp) }}</span>
						</template>
					</el-table-column>
					<el-table-column label="操作" width="230" align="center" fixed="right">
						<template #default="scope">
							<div class="list-button-class">
								<el-button size="small" type="text" @click="handleView(scope.row)" icon="View">查看</el-button>
								<el-button size="small" type="text" @click="handleEdit(scope.row)" icon="Edit">编辑</el-button>
								<el-button size="small" type="text" @click="handleDelete(scope.row)" icon="Delete">删除</el-button>
							</div>
						</template>
					</el-table-column>
				</el-table>
			</div>
		</el-card>

		<!-- yaml -->
		<el-drawer v-model="yamlDrawer" direction="rtl" :close-on-click-modal="false" :close-on-press-escape="false" size="56%">
			<template #header>
				<h4>Task 详情</h4>
			</template>
			<template #default>
				<div style="min-height: 300px; padding: 15px" v-loading="yamlDrawerLoading">
					<codeEditor v-if="!yamlDrawerLoading" v-model:modelValue="yamlContent" theme="ambiance-mobile"></codeEditor>
				</div>
			</template>
		</el-drawer>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { getTasks } from '/@/api/task';
import { parseTime } from '/@/utils/formatTime';
import codeEditor from '/@/components/codeEditor/index.vue';
// @ts-ignore
import yaml from 'js-yaml';

const listQuery = ref({
	name: '',
});
const list = ref<any>([]);
const yamlDrawer = ref(false);
const yamlDrawerLoading = ref(false);
const yamlContent = ref('');

const getList = () => {
	getTasks(listQuery.value).then((res) => {
		list.value = res.items;
	});
};

const handleView = (row: any) => {
	yamlContent.value = yaml.dump(row);
	yamlDrawer.value = true;
};
const handleEdit = (row: any) => {
	console.log(row);
};
const handleDelete = (row: any) => {
	console.log(row);
};

onMounted(() => {
	getList();
});
</script>

<style scoped lang="scss"></style>
