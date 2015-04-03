define([
'kernel',
'angular',
'./module'
],
function( core, ng, value){

    value
        .value('ErrorType', {
            DefaultError: 0,
            ExistsError: 1,
            ParamsError: 2,
            RequestError: 3,
            DbError: 4,
            EmptyError: 5
        })
});