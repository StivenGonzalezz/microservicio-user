Feature: Gestión de notificaciones

  Scenario: Crear una notificación exitosamente
    Given que preparo una petición válida para crear una notificación
    When llamo al servicio POST /notificaciones
    Then la respuesta debe tener código 201
    And el cuerpo debe coincidir con el schema crear_notificacion_schema.json

  Scenario: Consultar una notificación existente
    Given que existe una notificación creada previamente
    When llamo al servicio GET /notificaciones/{id}
    Then la respuesta debe tener código 200
    And el cuerpo debe coincidir con el schema consultar_notificacion_schema.json
