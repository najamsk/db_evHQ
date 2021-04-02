(function ($) {
  'use strict';

  $(document).ready(function () {
    $('#range').change(function () {
      console.log("changed")
      $('#start').toggle();
      $('#end').toggle();
    });

    $(".fa-chevron-up").hide();
    $(".fa-chevron-down").click(function () {
      $("#clientSearchBody").toggle(500);
      $(".fa-chevron-down").toggle();
      $(".fa-chevron-up").show();

    });
    $(".fa-chevron-up").click(function () {
      $(".fa-chevron-up").toggle();
      $("#clientSearchBody").toggle(500);
      $(".fa-chevron-down").toggle();

    });

    $("#clientCloseModal").click(function () {
      $("#create-app-modal").hide();

      $(".fa-chevron-down").click(function () {
        $("#clientSearchBody").show(500);
        $(".fa-chevron-down").toggle();
        $(".fa-chevron-up").show();

      });
      $(".fa-chevron-up").click(function () {
        $(".fa-chevron-up").toggle();
        $("#clientSearchBody").toggle(500);
        $(".fa-chevron-down").toggle();

      });

    });

    //datepickers

    $('#datetimepicker1, #datetimepicker12').datetimepicker();
    // $('.datetimepicker').datetimepicker({
    //   format: "YYYY-MM-DD hh:mm:ss A", //2019-08-28
    //   // timeZone: "Asia/Karachi"
    // });

    $('.datetimepicker').datetimepicker({
      //format: "YYYY-MM-DD hh:mm:ss A", 
    });

    $(".timepickerForm").submit(function () {
      console.log("intercept form grab all pickers and convert values to utc");
      var $this = $(this);
      var $pickers = $this.find(".datetimepicker").each(function (ele) {
        var picker = $(this)
        // console.log("about to print picker value");
        // console.log(picker.val());
        var pickerVal = picker.val();
        var cdate = new Date(pickerVal);

        var cdate2 = moment(pickerVal);
        var cdate3 = cdate2.valueOf();

        // console.log("cdate = " + cdate);
        var local = moment(cdate).valueOf();
        var utc = (moment(cdate).add(-(moment().utcOffset()), 'm'));
        utc = moment.parseZone(utc).utc().valueOf();
        // console.log('Local Timestamp is ' + local);
        // console.log('UTC Timestamp is ' + utc);
        picker.val(local);

      });
      // return false;

    });


    $('#datepicker-range, #datepicker-component, #datepicker-component2').datepicker();
    $('#datepicker-embeded').datepicker({
      daysOfWeekDisabled: '0,1'
    });
    $(function () {
      $('input[name="datetimes"]').daterangepicker({
        timePicker: true,
        startDate: moment().startOf('hour'),
        endDate: moment().startOf('hour').add(32, 'hour'),
        locale: {
          format: 'M/DD hh:mm A'
        }
      });
    });
    $('#daterangepicker, #daterangepicker20').daterangepicker({
      timePicker: true,
      timePickerIncrement: 30,
      format: 'MM/DD/YYYY h:mm A'
    }, function (start, end, label) {
      console.log(start.toISOString(), end.toISOString(), label);
    });
    $('#timepicker').timepicker().on('show.timepicker', function (e) {
      var widget = $('.bootstrap-timepicker-widget');
      widget.find('.glyphicon-chevron-up').removeClass().addClass('pg-arrow_maximize');
      widget.find('.glyphicon-chevron-down').removeClass().addClass('pg-arrow_minimize');
    });
    var nowTemp = new Date();
    var now = new Date(nowTemp.getFullYear(), nowTemp.getMonth(), nowTemp.getDate(), 0, 0, 0, 0);


    //   $('.datetimepicker').datetimepicker();

    //   $('.datepicker').datepicker({
    //     format: 'mm/dd/yyyy',
    //     startDate: '-3d'


    //     });



    //   $('#datetimepicker1, #datetimepicker11,  #daterangepicker2').datetimepicker();
    //   $(' .datepicker ').each(function() {
    //     $(this).datepicker('clearDates');
    //   });
    //   $('input[name="datetimes"]').daterangepicker({
    //         timePicker: true,
    //         startDate: moment().startOf('hour'),
    //         endDate: moment().startOf('hour').add(32, 'hour'),
    //         locale: {
    //         format: 'M/DD hh:mm A'
    //         }
    //     });

    // $(".list-view-wrapper").scrollbar();
    // $('[data-pages="search"]').search({
    //   searchField: '#overlay-search',
    //   closeButton: '.overlay-close',
    //   suggestions: '#overlay-suggestions',
    //   brand: '.brand',
    //   onSearchSubmit: function(searchString) {
    //     console.log("Search for: " + searchString);
    //   },
    //   onKeyEnter: function(searchString) {
    //     console.log("Live search for: " + searchString);
    //     var searchField = $('#overlay-search');
    //     var searchResults = $('.search-results');
    //     clearTimeout($.data(this, 'timer'));
    //     searchResults.fadeOut("fast");
    //     var wait = setTimeout(function() {
    //       searchResults.find('.result-name').each(function() {
    //         if (searchField.val().length != 0) {
    //           $(this).html(searchField.val());
    //           searchResults.fadeIn("fast");
    //         }
    //       });
    //     }, 500);
    //     $(this).data('timer', wait);
    //   }
    // });


  }); //doc ready


  $('.panel-collapse label').on('click', function (e) {
    e.stopPropagation();
  })


  //Modal Jquery


})(window.jQuery);
