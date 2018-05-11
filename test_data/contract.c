#include "ont.h"

int add(int a, int b ){
        return a + b;
}

char * concat(char * a, char * b){
	int lena = arrayLen(a);
	int lenb = arrayLen(b);
	char * res = (char *)malloc((lena + lenb)*sizeof(char));
	for (int i = 0 ;i < lena ;i++){
		res[i] = a[i];
	}

	for (int j = 0; j < lenb ;j++){
		res[lena + j] = b[j];
	}
	return res;
}


int sumArray(int * a, int * b){

	int res = 0;
	int lena = arrayLen(a);
	int lenb = arrayLen(b);

	for (int i = 0;i<lena;i++){
		res += a[i];
	}
	for (int j = 0;j<lenb;j++){
		res += b[j];
	}
	return res;
}


char * invoke(char * method,char * args){

    if (strcmp(method ,"init")==0 ){
            return "init success!";
    }

    if (strcmp(method, "add")==0){
        struct Params {
                int a;
                int b;
        };
        struct Params *p = (struct Params *)malloc(sizeof(struct Params));

        ONT_JsonUnmashalInput(p,sizeof(struct Params),args);
        int res = add(p->a,p->b);
        char * result = ONT_JsonMashalResult(res,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

	if(strcmp(method,"concat")==0){
		struct Params{
			char *a;
			char *b;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ONT_JsonUnmashalInput(p,sizeof(struct Params),args);
		char * res = concat(p->a,p->b);
		char * result = ONT_JsonMashalResult(res,"string",1);
		ONT_Runtime_Notify(result);
		return result;
	}
	
	if(strcmp(method,"sumArray")==0){
		struct Params{
			int *a;
			int *b;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ONT_JsonUnmashalInput(p,sizeof(struct Params),args);
		int res = sumArray(p->a,p->b);
		char * result = ONT_JsonMashalResult(res,"int",1);
		ONT_Runtime_Notify(result);
		return result;
	}

	if(strcmp(method,"addStorage")==0){

		struct Params{
			char * a;
			char * b;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ONT_JsonUnmashalInput(p,sizeof(struct Params),args);
		ONT_Storage_Put(p->a,p->b);
		char * result = ONT_JsonMashalResult("Done","string",1);
		ONT_Runtime_Notify(result);
		return result;
    }

	if(strcmp(method,"getStorage")==0){
		struct Params{
			char * a;
		};
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ONT_JsonUnmashalInput(p,sizeof(struct Params),args);
		char * value = ONT_Storage_Get(p->a);
		char * result = ONT_JsonMashalResult(value,"string",1);
		ONT_Runtime_Notify(result);
		return result;
	}

	if(strcmp(method,"deleteStorage")==0){

        struct Params{
                char * a;
        };
		struct Params *p = (struct Params *)malloc(sizeof(struct Params));
		ONT_JsonUnmashalInput(p,sizeof(struct Params),args);
        ONT_Storage_Delete(p->a);
        char * result = ONT_JsonMashalResult("Done","string",1);
        ONT_Runtime_Notify(result);
        return result;
    }
}