// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * @typedef {import('./models').main.Title} mainTitle
 * @typedef {import('./models').main.Person} mainPerson
 */

export const GreetService = {

		/**
		 * GreetService.Greet
		 * Greet does XYZ
         * @param name {string}
         * @param title {mainTitle}
		 * @returns {Promise<string>}
		 **/
	    Greet: function(name, title) { return wails.CallByName("main.GreetService.Greet", ...Array.prototype.slice.call(arguments, 0)); },

		/**
		 * GreetService.NewPerson
		 * NewPerson creates a new person
         * @param name {string}
		 * @returns {Promise<mainPerson>}
		 **/
	    NewPerson: function(name) { return wails.CallByName("main.GreetService.NewPerson", ...Array.prototype.slice.call(arguments, 0)); },
};

