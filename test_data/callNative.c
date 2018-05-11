#include "ont.h"
/*
*this is the common standard interface of ontology wasm contract
*/
char * invoke(char * method,char * args){

    if(strcmp(method,"transferont") == 0){
        char * from = ONT_ReadStringParam(args);
        char * to = ONT_ReadStringParam(args);
        long long ontAmount = ONT_ReadInt64Param(args);
        struct State{
            int ver;
            char * from;
            char * to ;
            long long amount;
        };

        struct Transfer {
            int ver;
            struct State * states
        };

        struct State * state = (struct State*)malloc(sizeof(struct State));
        state->ver = 1;
        state->amount = ontAmount;
        state->from = from;
        state->to = to;

        struct Transfer * transfer =(struct Transfer*)malloc(sizeof(struct Transfer));
        transfer->ver = 1;
        transfer->states = state;
        //1. call native contract to tranfer ont

        char * args = ONT_MarshalNativeParams(transfer);
        char * result = ONT_CallContract("ff00000000000000000000000000000000000001","","transfer",args);
        char * res;
        if (strcmp(result,"true")==0){
           res = ONT_JsonMashalResult("transfer succeed","string",1);
        }else{
           res = ONT_JsonMashalResult("transfer failed","string",0);
        }
        ONT_Runtime_Notify(res);
        return res;
    }
}