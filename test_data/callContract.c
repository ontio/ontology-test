#include "ont.h"
char * invoke(char * method,char * args){

    if (strcmp(method ,"init")==0 ){
            return "init success!";
    }

    if (strcmp(method, "getValue")==0){

        struct Input{
            char * key;
        };
        struct Input * input = (struct Input *)malloc(sizeof(struct Input));
        ONT_JsonUnmashalInput(input,sizeof(struct Input),args);

        struct Param{
            char * ptype;
            char * pvalue;
        };
        struct Param * newargs = (struct Param *)malloc(sizeof(struct Param));

        newargs -> ptype = "string";
        newargs -> pvalue = input->key;

        char * res = ONT_CallContract("90e45f98d7c49851fe7b6e33eef7a34b9507c493","","getStorage",ONT_JsonMashalParams(newargs));
        ONT_Runtime_Notify(res);
        return res;
    }
    if (strcmp(method,"putValue") == 0){

        struct Input{
            char * key;
            char * value;
        };
        struct Input * input = (struct Param *)malloc(sizeof(struct Input));
        ONT_JsonUnmashalInput(input,sizeof(struct Input),args);

        struct Param{
            char * ptype;
            char * pvalue;
        };

        struct Param * newargs = (struct Param *)malloc(sizeof(struct Param)*2);

        newargs[0].ptype = "string";
        newargs[0].pvalue = input->key;

        newargs[1].ptype = "string";
        newargs[1].pvalue = input->value;

        char * res = ONT_CallContract("90e45f98d7c49851fe7b6e33eef7a34b9507c493","","addStorage",ONT_JsonMashalParams(newargs));
        ONT_Runtime_Notify(res);
        return res;

    }
}