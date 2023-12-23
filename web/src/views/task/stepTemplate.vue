<template>
	<div class="main-container">
		<el-alert title="StepTemplate 是一种可重用的步骤定义模板，用于在流水线中定义和共享任务步骤。" class="mb10" type="info" show-icon />
		<el-card class="box-card">
			<div>
				<el-button type="primary" icon="ele-Plus" class="mr8">新建</el-button>
				<el-input v-model="listQuery.search_value" placeholder="请输入搜索项" @keyup.enter.native="getList" style="max-width: 300px">
					<template #append>
						<el-button icon="ele-Search" @click="getList"></el-button>
					</template>
				</el-input>
			</div>
			<div class="mt10">
				<el-table :data="list" border stripe style="width: 100%">
					<el-table-column prop="name" label="Date" />
					<el-table-column prop="icon" label="Name" />
				</el-table>
			</div>
		</el-card>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const listQuery = ref({});
const list = ref<any>([
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
	{ icon: 'ele-Plus', name: 'test1test1test1test1' },
]);

const getList = () => {};

const toLink = (item: any) => {
	console.log(item);
};
</script>

<style scoped lang="scss"></style>
