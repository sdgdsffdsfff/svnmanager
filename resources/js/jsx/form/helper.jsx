define(['kernel'], function(core){
    return {
        submitForm: $.noop,
        focusFirst: function(){
            var formBody = this.getRef('formBody'),
                fields = formBody.props.fields, field;

            if( fields && fields.length ){
                if( field = formBody.getRef(fields[0].name) ){
                    core.delay(function(){
                        field.focus();
                    }, 100)
                }
            }

            return this;
        },
        clearForm: function ( exclude ) {
            exclude = exclude || [];
            $.each(this.getRef('formBody').getRef(), function (k, v) {
                if( exclude.indexOf(k) > -1 ) return;

                if (v.props.value == '') {
                    v.clear();
                } else {
                    v.restore();
                }
            });
            return this;
        },
        setFormValue: function( formData ){
            var formBody = this.getRef('formBody');
            $.each( formData, function(key, value){
                var field;
                if( field = formBody.getRef(key) ){
                    field.setValue(value);
                }
            });
            return this;
        },
        getFormData: function () {
            var data = {}, defer = $.Deferred(), err = 0;
            $.each(this.getRef('formBody').getRef(), function (k, v) {
                data[k] = v.getValue();
                if (v.props.required && v.isEmpty()) {
                    err++;
                    v.focus();
                    defer.reject(k, v);
                    return false;
                }
            });
            if (err == 0) {
                defer.resolve(data)
            }
            return defer;
        }
    }
});